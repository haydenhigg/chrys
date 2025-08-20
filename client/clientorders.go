package client

import "time"

func (c *Client) PlaceOrder(config *OrderConfig, now time.Time) error {
	// normalize config
	if err := config.normalize(); err != nil {
		return err
	}

	// get latest price
	price, ok := c.Store.TryGetPrice(config.Pair, now)
	if !ok {
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
	case MARKET_BUY:
		quoteQuantity = config.Percent * balances[config.QuoteBalanceKey]
		baseQuantity = quoteQuantity / price
	case MARKET_SELL:
		baseQuantity = config.Percent * balances[config.BaseBalanceKey]
		quoteQuantity = baseQuantity * price
	}

	// place order
	if !config.IsDryRun {
		err = c.Connector.MarketOrder(string(config.Side), config.Pair, baseQuantity)
		if err != nil {
			return err
		}
	}

	// update balances
	switch config.Side {
	case MARKET_BUY:
		c.Store.Balances[config.QuoteBalanceKey] -= quoteQuantity
		c.Store.Balances[config.BaseBalanceKey] += baseQuantity * (1 - c.Fee)
	case MARKET_SELL:
		c.Store.Balances[config.BaseBalanceKey] -= baseQuantity
		c.Store.Balances[config.QuoteBalanceKey] += quoteQuantity * (1 - c.Fee)
	}

	return nil
}
