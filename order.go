package chrys

type OrderType string

const (
	BUY  OrderType = "buy"
	SELL           = "sell"
)

type Order struct {
	Pair    *Pair
	Percent float64
	IsLive  bool
	Type    OrderType
}

func NewOrder(pair *Pair, percent float64, isLive bool) *Order {
	return &Order{
		Pair:    pair,
		Percent: percent,
		IsLive:  isLive,
	}
}

func (order *Order) SetBuy() *Order {
	order.Type = BUY
	return order
}

func (order *Order) SetSell() *Order {
	order.Type = SELL
	return order
}

func (order *Order) normalize() {
	if order.Percent <= 0 {
		order.Percent = 0
	} else if order.Percent > 1 {
		order.Percent = 1
	}
}
