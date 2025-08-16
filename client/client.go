package client

import (
	"encoding/base64"
	"github.com/haydenhigg/clover/candle"
	"time"
)

type Client interface {
	GetBalances() (map[string]float64, error)
	GetCandlesSince(
		pair string,
		interval time.Duration,
		since time.Time,
	) ([]*candle.Candle, error)
	PlaceOrder(t string, pair string, quantity float64) error
}

type CachedClient struct {
	Client       Client
	BalanceStore map[string]float64
	CandleStore  map[string]map[time.Duration][]*candle.Candle
}

func (c *CachedClient) Balances() (map[string]float64, error) {
	if len(c.BalanceStore) == 0 {
		balances, err := c.Client.GetBalances()
		if err != nil {
			return nil, err
		}

		c.BalanceStore = balances
	}

	return c.BalanceStore, nil
}

func (c *CachedClient) Candles(
	pair string,
	interval time.Duration,
	n int,
) ([]*candle.Candle, error) {
	since := time.Now().Add(time.Duration(-n) * interval)
	return c.Client.GetCandlesSince(pair, interval, since)
}

// func (c *CachedClient) Order(type, pair string, quantity float64) error {
// 	return nil
// }

// client initialization
func NewKraken(key, secret string) (*CachedClient, error) {
	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}

	return &CachedClient{
		Client:       &Kraken{Key: []byte(key), Secret: decodedSecret},
		BalanceStore: map[string]float64{},
		CandleStore:  map[string]map[time.Duration][]*candle.Candle{},
	}, nil
}
