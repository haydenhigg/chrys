package chrys

import "fmt"

type Pair struct {
	Base  *Asset
	Quote *Asset
	Name  string
}

func NewPair(base, quote *Asset) *Pair {
	return &Pair{
		Base:  base,
		Quote: quote,
		Name: fmt.Sprintf("%s/%s", base.Symbol, quote.Symbol),
	}
}
