package chrys

import (
	"math"
	"time"
)

type Connector interface {
	FetchFramesSince(
		pair *Pair,
		interval time.Duration,
		since time.Time,
	) ([]*Frame, error)
	FetchBalances() (map[string]float64, error)
	PlaceMarketOrder(side, pair string, quantity float64) error
}

type Client struct {
	Connector Connector
	Frames    FrameCache
	Balances  map[string]float64
	Fee       float64
	IsLive    bool
}

func NewClient(connector Connector) *Client {
	return &Client{
		Connector: connector,
		Frames:    FrameCache{},
		Balances:  map[string]float64{},
	}
}

func (client *Client) SetFee(fee float64) *Client {
	client.Fee = fee
	return client
}

func (client *Client) SetIsLive(isLive bool) *Client {
	client.IsLive = isLive
	return client
}

// frames
func (client *Client) GetFramesSince(
	pair *Pair,
	interval time.Duration,
	t time.Time,
) ([]*Frame, error) {
	t = t.Truncate(interval)

	// check cache
	if frames, ok := client.Frames.GetSince(pair, interval, t); ok {
		return frames, nil
	}

	// retrieve from data source
	frames, err := client.Connector.FetchFramesSince(pair, interval, t)
	if err != nil {
		return nil, err
	}

	// cache retrieved data
	client.Frames.Set(pair, interval, frames)

	return frames, nil
}

func (client *Client) GetNFramesBefore(
	pair *Pair,
	interval time.Duration,
	n int,
	t time.Time,
) ([]*Frame, error) {
	t = t.Add(time.Duration(-n) * interval)

	frames, err := client.GetFramesSince(pair, interval, t)
	if err != nil {
		return nil, err
	}

	return frames[:n], nil
}

func (client *Client) GetPrice(pair *Pair, t time.Time) (float64, error) {
	price, ok := client.Frames.GetPriceAt(pair, t)
	if !ok {
		frames, err := client.GetNFramesBefore(pair, time.Minute, 1, t)
		if err != nil {
			return 0, err
		}

		price = frames[len(frames)-1].Close
	}

	return price, nil
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
	assets []*Asset,
	quote *Asset,
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

		var base *Asset = nil
		for _, asset := range assets {
			if asset.Code == baseCode {
				base = asset
				break
			}
		}

		if base == nil {
			continue
		}

		price, err := client.GetPrice(NewPair(base, quote), t)
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
	pair *Pair,
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
	price, err := client.GetPrice(pair, t)
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

func (client *Client) Buy(pair *Pair, percent float64, t time.Time) error {
	return client.Order(BUY, pair, percent, t)
}

func (client *Client) Sell(pair *Pair, percent float64, t time.Time) error {
	return client.Order(SELL, pair, percent, t)
}
