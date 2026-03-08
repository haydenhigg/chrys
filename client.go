package chrys

import (
	"github.com/haydenhigg/chrys/driver"
	"github.com/haydenhigg/chrys/store"
	"strings"
	"time"
)

type API interface {
	store.BalanceAPI
	store.FrameAPI
	MarketOrder(side, pair string, quantity float64) error
}

type Client struct {
	api      API
	Frames   *store.FrameStore
	Balances *store.BalanceStore
	Fee      float64
	IsLive   bool
}

// initializers
func NewClient(api API) *Client {
	return &Client{
		api:      api,
		Frames:   store.NewFrames(api),
		Balances: store.NewBalances(api),
	}
}

func NewKrakenClient(key, secret string) (*Client, error) {
	kraken, err := driver.NewKraken(key, secret)
	if err != nil {
		return nil, err
	}

	return NewClient(kraken).SetFee(0.004), nil
}

func NewHistoricalClient(dataRoot, nameFmt string) *Client {
	return NewClient(driver.NewHistorical(dataRoot, nameFmt))
}

// setters
func (client *Client) SetFee(fee float64) *Client {
	client.Fee = fee
	return client
}

func (client *Client) SetIsLive(isLive bool) *Client {
	client.IsLive = isLive
	return client
}

// methods
func (client *Client) Value(
	quoteAsset string,
	baseAssets []string,
	t time.Time,
) (float64, error) {
	values, err := client.Values(quoteAsset, baseAssets, t)
	if err != nil {
		return 0, err
	}

	total := 0.
	for _, value := range values {
		total += value
	}

	return total, nil
}

func (client *Client) Values(
	quoteAsset string,
	baseAssets []string,
	t time.Time,
) (map[string]float64, error) {
	values := make(map[string]float64, len(baseAssets))

	// check for assets
	if len(baseAssets) == 0 {
		return values, nil
	}

	// get all balances
	balances, err := client.Balances.Get()
	if err != nil {
		return values, err
	}

	// aggregate values
	for _, baseAsset := range baseAssets {
		balance, ok := balances[baseAsset]
		if !ok {
			continue
		}

		// add balance directly if asset is our quote asset
		if baseAsset == quoteAsset {
			values[baseAsset] = balance
			continue
		}

		price, err := client.Frames.GetPriceAt(baseAsset+"/"+quoteAsset, t)
		if err != nil {
			return values, err
		}

		values[baseAsset] = balance * price
	}

	return values, nil
}

type OrderSide string

const (
	BUY  OrderSide = "buy"
	SELL OrderSide = "sell"
)

func (client *Client) Order(
	side OrderSide,
	pair string,
	baseQuantity float64,
	t time.Time,
) error {
	// determine assets
	assets := strings.SplitN(pair, "/", 2)
	base, quote := assets[0], assets[1]

	// determine order quantities
	balances, err := client.Balances.Get()
	if err != nil {
		return err
	}

	baseQuantity = max(baseQuantity, 0)
	if side == SELL && baseQuantity > balances[base] {
		baseQuantity = balances[base]
	}

	price, err := client.Frames.GetPriceAt(pair, t)
	if err != nil {
		return err
	}

	quoteQuantity := baseQuantity * price
	if side == BUY && quoteQuantity > balances[quote] {
		quoteQuantity = balances[quote]
		baseQuantity = quoteQuantity / price
	}

	// place order
	if client.IsLive {
		err = client.api.MarketOrder(string(side), pair, baseQuantity)
		if err != nil {
			return err
		}
	}

	// update balances
	invFee := 1 - client.Fee
	switch side {
	case BUY:
		client.Balances.Set(map[string]float64{
			base:  baseQuantity * invFee,
			quote: -quoteQuantity,
		})
	case SELL:
		client.Balances.Set(map[string]float64{
			base:  -baseQuantity,
			quote: quoteQuantity * invFee,
		})
	}

	return nil
}

func (client *Client) Buy(pair string, quantity float64, t time.Time) error {
	return client.Order(BUY, pair, quantity, t)
}

func (client *Client) Sell(pair string, quantity float64, t time.Time) error {
	return client.Order(SELL, pair, quantity, t)
}

func (client *Client) OrderPct(
	side OrderSide,
	pair string,
	percent float64,
	t time.Time,
) error {
	balances, err := client.Balances.Get()
	if err != nil {
		return err
	}

	base := strings.SplitN(pair, "/", 2)[0]

	return client.Order(side, pair, percent*balances[base], t)
}

func (client *Client) BuyPct(pair string, percent float64, t time.Time) error {
	return client.OrderPct(BUY, pair, percent, t)
}

func (client *Client) SellPct(pair string, percent float64, t time.Time) error {
	return client.OrderPct(SELL, pair, percent, t)
}

// L1 normalize ReLU'd values to get weights from arbitrary values (softmax was
// avoided because it's not idempotent)
func scale(xs map[string]float64) map[string]float64 {
	sum := 0.
	for _, v := range xs {
		sum += max(v, 0) // ReLU
	}

	w := 0.
	if sum == 0 {
		sum = 1
		w = 1 / float64(len(xs))
	}

	scaledXs := make(map[string]float64, len(xs))
	for k, v := range xs {
		scaledXs[k] = max(v, 0)/sum + w
	}

	return scaledXs
}

func (client *Client) Reweight(
	quoteSymbol string,
	weights map[string]float64,
	t time.Time,
) error {
	// get current asset values
	symbols := make([]string, len(weights))
	i := 0
	for symbol := range weights {
		symbols[i] = symbol
		i++
	}

	values, err := client.Values(quoteSymbol, symbols, t)
	if err != nil {
		return err
	}

	totalValue := 0.
	for _, value := range values {
		totalValue += value
	}

	// calculate current and target asset weights
	targetWeights := scale(weights)
	currentWeights := scale(values)

	// order the difference
	for symbol, targetWeight := range targetWeights {
		// the quote asset cannot be bought or sold directly
		if symbol == quoteSymbol {
			continue
		}

		// calculate the difference
		weightDelta := targetWeight - currentWeights[symbol]
		if weightDelta == 0 {
			continue
		}

		pair := symbol + "/" + quoteSymbol
		price, err := client.Frames.GetPriceAt(pair, t)
		if err != nil {
			return err
		}

		quantityDelta := weightDelta * totalValue / price

		// place order
		if quantityDelta > 0 {
			err = client.Buy(pair, quantityDelta, t)
		} else if quantityDelta < 0 {
			err = client.Sell(pair, -quantityDelta, t)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
