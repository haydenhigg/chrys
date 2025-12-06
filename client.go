package chrys

import "time"

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
	Connector  Connector
	FrameCache map[string]map[time.Duration][]*Frame
	Balances   map[string]float64
	Fee        float64
}

func NewClient(connector Connector) *Client {
	return &Client{
		Connector:  connector,
		FrameCache: map[string]map[time.Duration][]*Frame{},
		Balances:   map[string]float64{},
	}
}

func (client *Client) SetFee(fee float64) *Client {
	client.Fee = fee
	return client
}

// frames
func (client *Client) getCachedFramesSince(
	pair *Pair,
	interval time.Duration,
	t time.Time,
) ([]*Frame, bool) {
	// ensure FrameCache[pair] exists
	if _, ok := client.FrameCache[pair.Name]; !ok {
		client.FrameCache[pair.Name] = map[time.Duration][]*Frame{}
		return nil, false
	}

	// assert that the frames exist and contain the time requested
	t = t.Truncate(interval)
	frames, ok := client.FrameCache[pair.Name][interval]
	if !ok || !frames[0].Time.Before(t.Add(interval)) {
		return nil, false
	}

	// chop off older frames
	for i, frame := range frames {
		if !frame.Time.Before(t) {
			return frames[i:], true
		}
	}

	return nil, false
}

func (client *Client) getCachedPriceAt(
	pair *Pair,
	t time.Time,
) (float64, bool) {
	// cycle through all cached intervals for an asset to see if any of them
	// have a price for the needed time to avoid searching and missing the
	// cache with a particular interval
	if intervalFrames, ok := client.FrameCache[pair.Name]; ok {
		frameTime := t.Truncate(time.Minute)

		for interval, frames := range intervalFrames {
			priorFrameTime := frameTime.Add(-interval)

			for _, frame := range frames {
				if frame.Time.Equal(priorFrameTime) {
					return frame.Close, true
				}
			}
		}
	} else {
		// ensure pair exists in cache
		client.FrameCache[pair.Name] = map[time.Duration][]*Frame{}
	}

	return 0, false
}

func (client *Client) GetFramesSince(
	pair *Pair,
	interval time.Duration,
	t time.Time,
) ([]*Frame, error) {
	t = t.Truncate(interval)

	// check cache
	if frames, ok := client.getCachedFramesSince(pair, interval, t); ok {
		return frames, nil
	}

	// retrieve from data source
	frames, err := client.Connector.FetchFramesSince(pair, interval, t)
	if err != nil {
		return nil, err
	}

	client.FrameCache[pair.Name][interval] = frames

	return frames, nil
}

func (client *Client) GetFrames(
	pair *Pair,
	interval time.Duration,
	t time.Time,
	n int,
) ([]*Frame, error) {
	t = t.Add(time.Duration(-n) * interval)

	frames, err := client.GetFramesSince(pair, interval, t)
	if err != nil {
		return nil, err
	}

	return frames[:n], nil
}

func (client *Client) GetPrice(pair *Pair, t time.Time) (float64, error) {
	price, ok := client.getCachedPriceAt(pair, t)
	if !ok {
		frames, err := client.GetFrames(pair, time.Minute, t, 1)
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
func (client *Client) PlaceOrder(order *Order, t time.Time) error {
	// get balances
	balances, err := client.GetBalances()
	if err != nil {
		return err
	}

	// get latest price
	price, err := client.GetPrice(order.Pair, t)
	if err != nil {
		return err
	}

	// determine quantities
	baseQuantity := order.Percent * balances[order.Pair.Base.Code]
	quoteQuantity := baseQuantity * price

	quoteBalance := balances[order.Pair.Quote.Code]
	if order.Type == MARKET_BUY && quoteQuantity > quoteBalance {
		// you can't spend more than you have
		quoteQuantity = quoteBalance
		baseQuantity = quoteQuantity / price
	}

	// place order
	if order.IsLive {
		err = client.Connector.PlaceMarketOrder(
			string(order.Type),
			order.Pair.Name,
			baseQuantity,
		)
		if err != nil {
			return err
		}
	}

	// update balances
	invFee := 1 - client.Fee

	switch order.Type {
	case MARKET_BUY:
		client.Balances[order.Pair.Quote.Code] -= quoteQuantity
		client.Balances[order.Pair.Base.Code] += baseQuantity * invFee
	case MARKET_SELL:
		client.Balances[order.Pair.Base.Code] -= baseQuantity
		client.Balances[order.Pair.Quote.Code] += quoteQuantity * invFee
	}

	return nil
}
