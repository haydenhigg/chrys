package chrys

import "fmt"

type Pair struct {
	base      string
	quote     string
	Symbol    string
	BaseCode  string
	QuoteCode string
}

func NewPair(base, quote string) *Pair {
	return &Pair{
		base:      base,
		quote:     quote,
		Symbol:    fmt.Sprintf("%s/%s", base, quote),
		BaseCode:  base,
		QuoteCode: quote,
	}
}

func (pair *Pair) Base() string {
	return pair.base
}

func (pair *Pair) Quote() string {
	return pair.base
}

func (pair *Pair) SetBaseCode(baseCode string) *Pair {
	pair.BaseCode = baseCode
	return pair
}

func (pair *Pair) SetQuoteCode(quoteCode string) *Pair {
	pair.QuoteCode = quoteCode
	return pair
}

func (pair *Pair) SetCodes(baseCode, quoteCode string) *Pair {
	return pair.SetBaseCode(baseCode).SetQuoteCode(quoteCode)
}
