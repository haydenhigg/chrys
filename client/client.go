package client

import "github.com/haydenhigg/chrys/client/driver"

type API interface {
	BalanceAPI
	FrameAPI
	MarketOrder(side, pair string, quantity float64) error
}

type Client struct {
	api      API
	Frames   *FrameStore
	Balances *BalanceStore
	Fee      float64
	IsLive   bool
}

// initializers
func New(api API) *Client {
	return &Client{
		api:      api,
		Frames:   NewFrameStore(api),
		Balances: NewBalanceStore(api),
	}
}

func NewKraken(key, secret string) (*Client, error) {
	kraken, err := driver.NewKraken(key, secret)
	if err != nil {
		return nil, err
	}

	return New(kraken).SetFee(0.004), nil
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
