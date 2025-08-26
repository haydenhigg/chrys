package client

import (
	"encoding/base64"
	"chrys"
	"chrys/client/connector"
	"time"
)

type Connector interface {
	FetchFramesSince(
		pair string,
		interval time.Duration,
		since time.Time,
	) ([]*chrys.Frame, error)
	FetchBalances() (map[string]float64, error)
	PlaceMarketOrder(side, pair string, quantity float64) error
}

type Client struct {
	Connector Connector
	Store     *Store
	Fee       float64
}

// initializers
func New(connector Connector, fee float64) *Client {
	return &Client{
		Connector: connector,
		Store: &Store{
			Frames:  map[string]map[time.Duration][]*chrys.Frame{},
			Balances: map[string]float64{},
		},
		Fee: fee,
	}
}

func NewKraken(key, secret string) (*Client, error) {
	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	connector := &connector.Kraken{Key: []byte(key), Secret: decodedSecret}

	return New(connector, 0.004), nil
}

func NewHistorical(dataRoot string, fee float64) *Client {
	return New(&connector.Historical{DataRoot: dataRoot}, fee)
}

// frames
func (c *Client) GetFramesSince(
	pair string,
	interval time.Duration,
	since time.Time,
) ([]*chrys.Frame, error) {
	since = since.Truncate(interval)

	// check store
	if frames, ok := c.Store.TryGetFramesSince(pair, interval, since); ok {
		return frames, nil
	}

	// retrieve from data source
	frames, err := c.Connector.FetchFramesSince(pair, interval, since)
	if err != nil {
		return nil, err
	}

	c.Store.Frames[pair][interval] = frames

	return frames, nil
}

func (c *Client) GetFrames(
	pair string,
	interval time.Duration,
	now time.Time,
	n int,
) ([]*chrys.Frame, error) {
	since := now.Add(time.Duration(-n) * interval)
	frames, err := c.GetFramesSince(pair, interval, since)
	if err != nil {
		return nil, err
	}

	return frames[:n], nil
}

// balances
func (c *Client) GetBalances() (map[string]float64, error) {
	// check store
	if balances, ok := c.Store.TryGetBalances(); ok {
		return balances, nil
	}

	// retrieve from data source
	balances, err := c.Connector.FetchBalances()
	if err != nil {
		return nil, err
	}

	c.Store.Balances = balances

	return balances, nil
}

// ordering
func (c *Client) PlaceOrder(config *OrderConfig, now time.Time) error {
	// normalize config
	if err := config.normalize(); err != nil {
		return err
	}

	// get latest price
	price, ok := c.Store.TryGetPriceAt(config.Pair, now)
	if !ok {
		frames, err := c.GetFrames(config.Pair, time.Minute, now, 1)
		if err != nil {
			return err
		}

		price = frames[len(frames)-1].Close
	}

	// determine quantities
	balances, err := c.GetBalances()
	if err != nil {
		return err
	}

	var baseQuantity, quoteQuantity float64

	switch config.Side {
	case MARKET_BUY:
		quoteQuantity = config.Percent * balances[config.QuoteBalanceKey]
		baseQuantity = quoteQuantity / price
	case MARKET_SELL:
		baseQuantity = config.Percent * balances[config.BaseBalanceKey]
		quoteQuantity = baseQuantity * price
	}

	// place order
	if !config.IsDryRun {
		err = c.Connector.PlaceMarketOrder(string(config.Side), config.Pair, baseQuantity)
		if err != nil {
			return err
		}
	}

	// update balances
	switch config.Side {
	case MARKET_BUY:
		c.Store.Balances[config.QuoteBalanceKey] -= quoteQuantity
		c.Store.Balances[config.BaseBalanceKey] += baseQuantity * (1 - c.Fee)
	case MARKET_SELL:
		c.Store.Balances[config.BaseBalanceKey] -= baseQuantity
		c.Store.Balances[config.QuoteBalanceKey] += quoteQuantity * (1 - c.Fee)
	}

	return nil
}
