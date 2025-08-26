package chrys

import "fmt"

type Pair struct {
	Base         string
	Quote        string
	BalanceBase  string
	BalanceQuote string
}

func NewPair(base, quote string) *Pair {
	return &Pair{
		Base:         base,
		Quote:        quote,
		BalanceBase:  base,
		BalanceQuote: quote,
	}
}

func (pair *Pair) SetPair(base, quote string) *Pair {
	pair.Base = base
	pair.Quote = quote
	return pair
}

func (pair *Pair) SetBalancePair(base, quote string) *Pair {
	pair.BalanceBase = base
	pair.BalanceQuote = quote
	return pair
}

func (pair *Pair) String() string {
	return fmt.Sprintf("%s/%s", pair.Base, pair.Quote)
}
