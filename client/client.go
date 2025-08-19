package client

import (
	"encoding/base64"
	"github.com/haydenhigg/clover/candle"
	"github.com/haydenhigg/clover/client/connector"
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
	Fee       float64
}

func New(connector Connector) *Client {
	return &Client{
		Connector: connector,
		Store: &Store{
			Candles:  map[string]map[time.Duration][]*candle.Candle{},
			Balances: map[string]float64{},
		},
	}
}

func NewKraken(key, secret string) (*Client, error) {
	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	client := New(&connector.Kraken{Key: []byte(key), Secret: decodedSecret})
	client.Fee = 0.004

	return client, nil
}
