package client

import "github.com/haydenhigg/chrys"

type OrderSide string

const (
	MARKET_BUY  OrderSide = "buy"
	MARKET_SELL           = "sell"
)

type OrderConfig struct {
	Side     OrderSide
	Pair     *chrys.Pair
	Percent  float64
	IsDryRun bool
}

func (config *OrderConfig) validatePercent() error {
	if config.Percent <= 0 {
		config.Percent = 0
	} else if config.Percent > 1 {
		config.Percent = 1
	}

	return nil
}

func (config *OrderConfig) normalize() error {
	if err := config.validatePercent(); err != nil {
		return err
	}

	return nil
}
