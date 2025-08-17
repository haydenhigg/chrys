package client

import (
	"github.com/haydenhigg/clover/candle"
	"github.com/haydenhigg/clover/client/connector"
	"encoding/base64"
	"time"
)

type Connector interface {
	GetCandlesSince(
		pair string,
		interval time.Duration,
		since time.Time,
	) ([]*candle.Candle, error)
	GetBalances() (map[string]float64, error)
	MarketOrder(side, pair string, quantity float64) error
}

type Client struct {
	Connector Connector
	Store     *Store
}

func New(connector Connector) *Client {
	return &Client{
		Connector: connector,
		Store: &Store{
			Candles: map[string]map[time.Duration][]*candle.Candle{},
			Balances: map[string]float64{},
		},
	}
}

func NewKraken(key, secret string) (*Client, error) {
	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	return New(&connector.Kraken{Key: []byte(key), Secret: decodedSecret}), nil
}

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

func (c *Client) GetBalances() (map[string]float64, error) {
	// check store
	if balances, ok := c.Store.TryGetBalances(); ok {
		return balances, nil
	}

	// retrieve from data source
	balances, err := c.Connector.GetBalances()
	if err != nil {
		return nil, err
	}

	c.Store.Balances = balances

	return balances, nil
}

func (c *Client) MarketOrder(side, pair string, quantity float64) error {
	return c.Connector.MarketOrder(side, pair, quantity)
}
