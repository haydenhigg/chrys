package store

type BalanceAPI interface {
	FetchBalances() (map[string]float64, error)
}

type BalanceStore struct {
	api      BalanceAPI
	balances map[string]float64
	aliases  map[string]string
}

func NewBalances(api BalanceAPI) *BalanceStore {
	return &BalanceStore{
		api:      api,
		balances: map[string]float64{},
		aliases:  map[string]string{},
	}
}

func (store *BalanceStore) Get() (map[string]float64, error) {
	// check if balances is not empty
	if len(store.balances) > 0 {
		return store.balances, nil
	}

	// retrieve from data source
	balances, err := store.api.FetchBalances()
	if err != nil {
		panic(err)
	}

	// cache retrieved data
	store.Set(balances)

	return balances, nil
}

func (store *BalanceStore) Set(balances map[string]float64) *BalanceStore {
	for asset, balance := range balances {
		store.balances[asset] += balance

		if alias, ok := store.aliases[asset]; ok {
			store.balances[alias] += balance
		}
	}

	return store
}

func (store *BalanceStore) Alias(asset, assetAlias string) *BalanceStore {
	if asset != assetAlias {
		store.aliases[asset] = assetAlias
		store.aliases[assetAlias] = asset
	}

	return store
}

// func (store *BalanceStore) TotalValue(quote string, t time.Time) (float64, error) {
// 	balances, err := store.Get()
// 	if err != nil {
// 		return 0, err
// 	}

// 	total := 0.

// 	for baseCode, balance := range balances {
// 		if baseCode == quote.Code {
// 			total += balance
// 			continue
// 		}

// 		var base *chrys.Asset = nil
// 		for _, asset := range assets {
// 			if asset.Code == baseCode {
// 				base = asset
// 				break
// 			}
// 		}

// 		if base == nil {
// 			continue
// 		}

// 		price, err := store.GetPriceAt(base.Symbol + "/" + quote.Symbol, t)
// 		if err != nil {
// 			return 0, err
// 		}

// 		total += balance * price
// 	}

// 	return total, nil
// }
