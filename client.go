package chrys

import (
	"github.com/haydenhigg/chrys/driver"
	"github.com/haydenhigg/chrys/store"
	"strings"
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

func NewHistoricalClient(dataRoot, nameFmt string) *Client {
	return NewClient(driver.NewHistorical(dataRoot, nameFmt))
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
func (client *Client) Value(assets []string, t time.Time) (float64, error) {
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

type OrderSide string

const (
	BUY  OrderSide = "buy"
	SELL OrderSide = "sell"
)

func (client *Client) Order(
	side OrderSide,
	pair string,
	percent float64,
	t time.Time,
) error {
	// determine assets
	assets := strings.SplitN(pair, "/", 2)
	base, quote := assets[0], assets[1]

	// determine order quantities
	adjPercent := max(percent, 0)
	if side == SELL {
		adjPercent = min(percent, 1)
	}

	balances, err := client.Balances.Get()
	if err != nil {
		return err
	}

	price, err := client.Frames.GetPriceAt(pair, t)
	if err != nil {
		return err
	}

	baseQuantity := adjPercent * balances[base]
	quoteQuantity := baseQuantity * price

	if side == BUY && quoteQuantity > balances[quote] {
		// you can't spend more than you have
		quoteQuantity = balances[quote]
		baseQuantity = quoteQuantity / price
	}

	// place order
	if client.IsLive {
		err = client.api.MarketOrder(string(side), pair, baseQuantity)
		if err != nil {
			return err
		}
	}

	// update balances
	invFee := 1 - client.Fee

	switch side {
	case BUY:
		client.Balances.Set(map[string]float64{
			base:  baseQuantity * invFee,
			quote: -quoteQuantity,
		})
	case SELL:
		client.Balances.Set(map[string]float64{
			base:  -baseQuantity,
			quote: quoteQuantity * invFee,
		})
	}

	return nil
}

func (client *Client) Buy(pair string, percent float64, t time.Time) error {
	return client.Order(BUY, pair, percent, t)
}

func (client *Client) Sell(pair string, percent float64, t time.Time) error {
	return client.Order(SELL, pair, percent, t)
}
