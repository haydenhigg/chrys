package client

import (
	"github.com/haydenhigg/clover/candle"
	"time"
)

func (c *Client) GetCandlesSince(
	pair string,
	interval time.Duration,
	since time.Time,
) ([]*candle.Candle, error) {
	since = since.Truncate(interval)

	// check store
	if candles, ok := c.Store.TryGetCandlesSince(pair, interval, since); ok {
		return candles, nil
	}

	// retrieve from data source
	candles, err := c.Connector.GetCandlesSince(pair, interval, since)
	if err != nil {
		return nil, err
	}

	c.Store.Candles[pair][interval] = candles

	return candles, nil
}

func (c *Client) GetCandles(
	pair string,
	interval time.Duration,
	now time.Time,
	n int,
) ([]*candle.Candle, error) {
	since := now.Add(time.Duration(-n) * interval)
	return c.GetCandlesSince(pair, interval, since)
}
