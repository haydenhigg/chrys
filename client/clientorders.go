package client

import (
	"errors"
	"strings"
	"time"
)

type OrderConfig struct {
	Side     string
	Pair     string
	Percent  float64
	IsDryRun bool

	BaseCode  string // for balances
	QuoteCode string // for balances
}

func (c *Client) MarketOrder(config *OrderConfig) error {
	if config.Side == "" || config.Percent == 0 {
		return nil
	}

	// normalize config
	if config.Percent < 0 || config.Percent > 1 {
		return errors.New("order percent must be in range [0, 1]")
	} else if config.Pair == "" {
		return errors.New("order pair must be defined")
	}

	symbols := strings.SplitN(config.Pair, "/", 2)
	if len(symbols) != 2 {
		return errors.New("order pair must follow format \"AAA/BBB\"")
	}

	if config.BaseCode == "" {
		config.BaseCode = symbols[0]
	}
	if config.QuoteCode == "" {
		config.QuoteCode = symbols[1]
	}

	// get latest price
	now := time.Now()
	var price float64

	if intervalCandles, ok := c.Store.Candles[config.Pair]; ok {
		for interval, candles := range intervalCandles {
			latestCandle := candles[len(candles)-1]
			if latestCandle.Time == now.Truncate(interval).Add(-interval) {
				price = latestCandle.Close
			}
		}
	}

	if price == 0 {
		candles, err := c.GetCandles(config.Pair, time.Minute, now, 1)
		if err != nil {
			return err
		}

		price = candles[len(candles)-1].Close
	}

	// determine quantities
	balances, err := c.GetBalances()
	if err != nil {
		return err
	}

	var baseQuantity, quoteQuantity float64

	switch config.Side {
	case "buy":
		quoteQuantity = config.Percent * balances[config.QuoteCode]
		baseQuantity = quoteQuantity / price
	case "sell":
		baseQuantity = config.Percent * balances[config.BaseCode]
		quoteQuantity = baseQuantity * price
	default:
		return errors.New("order side must be either \"buy\" or \"sell\"")
	}

	// place order
	if !config.IsDryRun {
		err = c.Connector.MarketOrder(config.Side, config.Pair, baseQuantity)
		if err != nil {
			return err
		}
	}

	// update balances
	switch config.Side {
	case "buy":
		c.Store.Balances[config.QuoteCode] -= quoteQuantity
		c.Store.Balances[config.BaseCode] += baseQuantity * (1 - c.Fee)
	case "sell":
		c.Store.Balances[config.BaseCode] -= baseQuantity
		c.Store.Balances[config.QuoteCode] += quoteQuantity * (1 - c.Fee)
	}

	return nil
}
