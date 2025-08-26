package chrys

import "fmt"

type Pair struct {
	Base      string
	Quote     string
	BaseCode  string
	QuoteCode string
}

func NewPair(base, quote string) *Pair {
	return &Pair{
		Base:      base,
		Quote:     quote,
		BaseCode:  base,
		QuoteCode: quote,
	}
}

func (pair *Pair) SetBase(base string) *Pair {
	pair.Base = base
	return pair
}

func (pair *Pair) SetQuote(quote string) *Pair {
	pair.Quote = quote
	return pair
}

func (pair *Pair) Set(base, quote string) *Pair {
	return pair.SetBase(base).SetQuote(quote)
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

func (pair *Pair) String() string {
	return fmt.Sprintf("%s/%s", pair.Base, pair.Quote)
}
