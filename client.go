package chrys

import (
	"fmt"
	"github.com/haydenhigg/chrys/driver"
	"github.com/haydenhigg/chrys/store"
	"time"
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

// methods
func (client *Client) TotalValue(
	assets []string,
	quote string,
	t time.Time,
) (float64, error) {
	balances, err := client.Balances.Get()
	if err != nil {
		return 0, err
	}

	total := 0.
	for _, base := range assets {
		fmt.Println(base, quote)

		// get balance for asset
		balance, ok := balances[base]
		if !ok {
			continue
		}

		fmt.Println(base, quote, balance)

		// add balance directly if asset is quote asset
		if base == quote {
			total += balance
			continue
		}

		// add balance for asset given pair price
		price, err := client.Frames.GetPriceAt(base+"/"+quote, t)
		if err != nil {
			return total, err
		}

		fmt.Println(base, quote, balance*price)

		total += balance * price
	}

	return total, nil
}
