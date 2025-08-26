package client

import (
	"encoding/base64"
	"github.com/haydenhigg/chrys"
	"github.com/haydenhigg/chrys/client/connector"
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
			Frames:   map[string]map[time.Duration][]*chrys.Frame{},
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
	series *chrys.Series,
	since time.Time,
) ([]*chrys.Frame, error) {
	since = since.Truncate(series.Interval)

	// check store
	if frames, ok := c.Store.TryGetFramesSince(series, since); ok {
		return frames, nil
	}

	// retrieve from data source
	frames, err := c.Connector.FetchFramesSince(
		series.Pair.String(),
		series.Interval,
		since,
	)
	if err != nil {
		return nil, err
	}

	c.Store.Frames[series.Pair.String()][series.Interval] = frames

	return frames, nil
}

func (c *Client) GetFrames(
	series *chrys.Series,
	now time.Time,
	n int,
) ([]*chrys.Frame, error) {
	since := now.Add(time.Duration(-n) * series.Interval)
	frames, err := c.GetFramesSince(series, since)
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
	price, ok := c.Store.TryGetPriceAt(config.Pair.String(), now)
	if !ok {
		series := chrys.NewSeries(config.Pair.String(), time.Minute)
		frames, err := c.GetFrames(series, now, 1)
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
		quoteQuantity = config.Percent * balances[config.Pair.QuoteCode]
		baseQuantity = quoteQuantity / price
	case MARKET_SELL:
		baseQuantity = config.Percent * balances[config.Pair.BaseCode]
		quoteQuantity = baseQuantity * price
	}

	// place order
	if !config.IsDryRun {
		err = c.Connector.PlaceMarketOrder(string(config.Side), config.Pair.String(), baseQuantity)
		if err != nil {
			return err
		}
	}

	// update balances
	invFee := 1 - c.Fee

	switch config.Side {
	case MARKET_BUY:
		c.Store.Balances[config.Pair.QuoteCode] -= quoteQuantity
		c.Store.Balances[config.Pair.BaseCode] += baseQuantity * invFee
	case MARKET_SELL:
		c.Store.Balances[config.Pair.BaseCode] -= baseQuantity
		c.Store.Balances[config.Pair.QuoteCode] += quoteQuantity * invFee
	}

	return nil
}
