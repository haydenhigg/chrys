package chrys

import "time"

type Connector interface {
	FetchFramesSince(
		series *Series,
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
	series *Series,
	t time.Time,
) ([]*Frame, bool) {
	// ensure FrameCache[pair] exists
	if _, ok := client.FrameCache[series.Pair.Name]; !ok {
		client.FrameCache[series.Pair.Name] = map[time.Duration][]*Frame{}
		return nil, false
	}

	// assert that the frames exist and contain the time requested
	t = t.Truncate(series.Interval)
	frames, ok := client.FrameCache[series.Pair.Name][series.Interval]
	if !ok || !frames[0].Time.Before(t.Add(series.Interval)) {
		return nil, false
	}

	// chop off oldest frames
	for i, frame := range frames {
		if !frame.Time.Before(t) {
			return frames[i:], true
		}
	}

	return nil, false
}

func (client *Client) GetFramesSince(
	series *Series,
	t time.Time,
) ([]*Frame, error) {
	t = t.Truncate(series.Interval)

	// check cache
	if frames, ok := client.getCachedFramesSince(series, t); ok {
		return frames, nil
	}

	// retrieve from data source
	frames, err := client.Connector.FetchFramesSince(series, t)
	if err != nil {
		return nil, err
	}

	client.FrameCache[series.Pair.Name][series.Interval] = frames

	return frames, nil
}

func (client *Client) GetFrames(
	series *Series,
	t time.Time,
	n int,
) ([]*Frame, error) {
	t = t.Add(time.Duration(-n) * series.Interval)

	frames, err := client.GetFramesSince(series, t)
	if err != nil {
		return nil, err
	}

	return frames[:n], nil
}

// balances
func (client *Client) getStoredBalances() (map[string]float64, bool) {
	if client.Balances == nil || len(client.Balances) == 0 {
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

// ordering
func (client *Client) getCachedPriceAt(
	pair *Pair,
	t time.Time,
) (float64, bool) {
	if intervalFrames, ok := client.FrameCache[pair.Name]; ok {
		// scan all intervals for pair
		for interval, frames := range intervalFrames {
			// look for the frame before (because using the Close)
			lastFrameTime := t.Truncate(interval).Add(-interval)

			for _, frame := range frames {
				if frame.Time.Equal(lastFrameTime) {
					return frame.Close, true
				}
			}
		}
	} else {
		// ensure FrameCache[pair] exists
		client.FrameCache[pair.Name] = map[time.Duration][]*Frame{}
	}

	return 0, false
}

func (client *Client) PlaceOrder(order *Order, t time.Time) error {
	order.normalize()

	// get latest price
	price, ok := client.getCachedPriceAt(order.Pair, t)
	if !ok {
		minSeries := NewSeries(order.Pair, time.Minute)
		frames, err := client.GetFrames(minSeries, t, 1)
		if err != nil {
			return err
		}

		price = frames[len(frames)-1].Close
	}

	// determine quantities
	balances, err := client.GetBalances()
	if err != nil {
		return err
	}

	var baseQuantity, quoteQuantity float64

	switch order.Type {
	case BUY:
		quoteQuantity = order.Percent * balances[order.Pair.Quote.Code]
		baseQuantity = quoteQuantity / price
	case SELL:
		baseQuantity = order.Percent * balances[order.Pair.Base.Code]
		quoteQuantity = baseQuantity * price
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
	case BUY:
		client.Balances[order.Pair.Quote.Code] -= quoteQuantity
		client.Balances[order.Pair.Base.Code] += baseQuantity * invFee
	case SELL:
		client.Balances[order.Pair.Base.Code] -= baseQuantity
		client.Balances[order.Pair.Quote.Code] += quoteQuantity * invFee
	}

	return nil
}
