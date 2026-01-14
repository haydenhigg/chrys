package client

import (
	"github.com/haydenhigg/chrys"
	"github.com/haydenhigg/chrys/client/connector"
	"math"
	"time"
)

type Connector interface {
	FetchFramesSince(
		pair *chrys.Pair,
		interval time.Duration,
		since time.Time,
	) ([]*chrys.Frame, error)
	FetchBalances() (map[string]float64, error)
	PlaceMarketOrder(side, pair string, quantity float64) error
}

type Client struct {
	Connector Connector
	Frames    *Frames
	Balances  map[string]float64
	Fee       float64
	IsLive    bool
}

// initializers
func New(connector Connector) *Client {
	return &Client{
		Connector: connector,
		Frames:    NewFrames(connector),
		Balances:  map[string]float64{},
	}
}

func NewKraken(key, secret string) (*Client, error) {
	kraken, err := connector.NewKraken(key, secret)
	if err != nil {
		return nil, err
	}

	return New(kraken), nil
}

// setters
func (client *Client) SetFee(fee float64) *Client {
	client.Fee = fee
	return client
}

func (client *Client) SetIsLive(isLive bool) *Client {
	client.IsLive = isLive
	return client
}

// balances
func (client *Client) getStoredBalances() (map[string]float64, bool) {
	if len(client.Balances) == 0 {
		return nil, false
	}

	return client.Balances, true
}

func (client *Client) GetBalances() (map[string]float64, error) {
	// check store
	if balances, ok := client.getStoredBalances(); ok {
		return balances, nil
	}

	// retrieve from data source
	balances, err := client.Connector.FetchBalances()
	if err != nil {
		return nil, err
	}

	client.Balances = balances

	return balances, nil
}

func (client *Client) GetTotalValue(
	assets []*chrys.Asset,
	quote *chrys.Asset,
	t time.Time,
) (float64, error) {
	balances, err := client.GetBalances()
	if err != nil {
		return 0, err
	}

	total := 0.

	for baseCode, balance := range balances {
		if baseCode == quote.Code {
			total += balance
			continue
		}

		var base *chrys.Asset = nil
		for _, asset := range assets {
			if asset.Code == baseCode {
				base = asset
				break
			}
		}

		if base == nil {
			continue
		}

		price, err := client.Frames.GetPriceAt(chrys.NewPair(base, quote), t)
		if err != nil {
			return 0, err
		}

		total += balance * price
	}

	return total, nil
}

// ordering
type OrderSide string

const (
	BUY  OrderSide = "buy"
	SELL OrderSide = "sell"
)

func (client *Client) Order(
	side OrderSide,
	pair *chrys.Pair,
	percent float64,
	t time.Time,
) error {
	// clamp percent to [0, 1]
	percent = math.Min(math.Max(percent, 0), 1)

	// get balances
	balances, err := client.GetBalances()
	if err != nil {
		return err
	}

	// get latest price
	price, err := client.Frames.GetPriceAt(pair, t)
	if err != nil {
		return err
	}

	// determine quantities
	baseQuantity := percent * balances[pair.Base.Code]
	quoteQuantity := baseQuantity * price

	quoteBalance := balances[pair.Quote.Code]
	if side == BUY && quoteQuantity > quoteBalance {
		// you can't spend more than you have
		quoteQuantity = quoteBalance
		baseQuantity = quoteQuantity / price
	}

	// place order
	if client.IsLive {
		err = client.Connector.PlaceMarketOrder(
			string(side),
			pair.Name,
			baseQuantity,
		)
		if err != nil {
			return err
		}
	}

	// update balances
	invFee := 1 - client.Fee

	switch side {
	case BUY:
		client.Balances[pair.Quote.Code] -= quoteQuantity
		client.Balances[pair.Base.Code] += baseQuantity * invFee
	case SELL:
		client.Balances[pair.Base.Code] -= baseQuantity
		client.Balances[pair.Quote.Code] += quoteQuantity * invFee
	}

	return nil
}

func (client *Client) Buy(
	pair *chrys.Pair,
	percent float64,
	t time.Time,
) error {
	return client.Order(BUY, pair, percent, t)
}

func (client *Client) Sell(
	pair *chrys.Pair,
	percent float64,
	t time.Time,
) error {
	return client.Order(SELL, pair, percent, t)
}
