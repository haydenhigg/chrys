package chrys

type Asset struct {
	Symbol string
	Code   string // for trading
}

func NewAsset(symbol, code string) *Asset {
	return &Asset{
		Symbol: symbol,
		Code: code,
	}
}
