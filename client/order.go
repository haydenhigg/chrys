package client

import (
	"math"
	"strings"
	"time"
)

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
	// clamp percent to [0, 1]
	percent = math.Min(math.Max(percent, 0), 1)

	// determine assets
	assets := strings.SplitN(pair, "/", 2)
	base, quote := assets[0], assets[1]

	// determine order quantities
	balances, err := client.Balances.Get()
	if err != nil {
		return err
	}

	price, err := client.Frames.GetPriceAt(pair, t)
	if err != nil {
		return err
	}

	baseQuantity := percent * balances[base]
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
			base:  balances[base] + baseQuantity*invFee,
			quote: balances[quote] - quoteQuantity,
		})
	case SELL:
		client.Balances.Set(map[string]float64{
			base:  balances[base] - baseQuantity,
			quote: balances[quote] + quoteQuantity*invFee,
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
