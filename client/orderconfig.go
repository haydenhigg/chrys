package client

import (
	"errors"
	"strings"
)

type OrderSide string

const (
	MARKET_BUY  OrderSide = "buy"
	MARKET_SELL           = "sell"
)

type OrderConfig struct {
	Side            OrderSide
	Pair            string
	BaseBalanceKey  string // for balances
	QuoteBalanceKey string // for balances
	Percent         float64
	IsDryRun        bool
}

func (config *OrderConfig) validatePair() ([2]string, error) {
	symbols := strings.Split(config.Pair, "/")
	if len(symbols) != 2 {
		return [2]string{}, errors.New("order pair must follow format 'BASE/QUOTE'")
	}

	return [2]string{symbols[0], symbols[1]}, nil
}

func (config *OrderConfig) validatePercent() error {
	if config.Percent <= 0 || config.Percent > 1 {
		return errors.New("order percent must be in range (0, 1]")
	}

	return nil
}

func (config *OrderConfig) normalize() error {
	if symbols, err := config.validatePair(); err != nil {
		return err
	} else {
		if config.BaseBalanceKey == "" {
			config.BaseBalanceKey = symbols[0]
		}

		if config.QuoteBalanceKey == "" {
			config.QuoteBalanceKey = symbols[1]
		}
	}

	if err := config.validatePercent(); err != nil {
		return err
	}

	return nil
}
