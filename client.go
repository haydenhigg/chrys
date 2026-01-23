package chrys

import (
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

func NewHistoricalClient(dataRoot, nameFmt string) (*Client, error) {
	return NewClient(driver.NewHistorical(dataRoot, nameFmt)), nil
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
	t time.Time,
) (float64, error) {
	// check for assets
	if len(assets) == 0 {
		return 0, nil
	}

	// set quote asset
	quote := assets[0]

	// get all balances
	balances, err := client.Balances.Get()
	if err != nil {
		return 0, err
	}

	// sum balances
	total := 0.
	for _, base := range assets {
		// get balance for asset
		balance, ok := balances[base]
		if !ok {
			continue
		}

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

		total += balance * price
	}

	return total, nil
}
