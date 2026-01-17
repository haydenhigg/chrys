package chrys

import (
	"github.com/haydenhigg/chrys/driver"
	"github.com/haydenhigg/chrys/store"
)

type API interface {
	store.BalanceAPI
	store.FrameAPI
	MarketOrder(side, pair string, quantity float64) error
}

type Client struct {
	api      API
	Frames   *store.FrameStore
	Balances *store.BalanceStore
	Fee      float64
	IsLive   bool
}

// initializers
func NewClient(api API) *Client {
	return &Client{
		api:      api,
		Frames:   store.NewFrames(api),
		Balances: store.NewBalances(api),
	}
}

func NewKrakenClient(key, secret string) (*Client, error) {
	kraken, err := driver.NewKraken(key, secret)
	if err != nil {
		return nil, err
	}

	return NewClient(kraken).SetFee(0.004), nil
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
