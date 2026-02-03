package store

type BalanceAPI interface {
	FetchBalances() (map[string]float64, error)
}

type BalanceStore struct {
	api      BalanceAPI
	Balances map[string]float64
	Aliases  map[string]string
}

func NewBalances(api BalanceAPI) *BalanceStore {
	return &BalanceStore{
		api:      api,
		Balances: map[string]float64{},
		Aliases:  map[string]string{},
	}
}

func (store *BalanceStore) Get() (map[string]float64, error) {
	// check if balances is not empty
	if len(store.Balances) > 0 {
		return store.Balances, nil
	}

	// retrieve from data source
	balances, err := store.api.FetchBalances()
	if err != nil {
		panic(err)
	}

	// cache retrieved data
	store.Set(balances)

	return store.Balances, nil
}

func (store *BalanceStore) Set(balances map[string]float64) *BalanceStore {
	// update all balances additively
	for asset, balance := range balances {
		store.Balances[asset] += balance

		// update alias balances
		if alias, ok := store.Aliases[asset]; ok {
			store.Balances[alias] += balance
		}
	}

	return store
}

func (store *BalanceStore) Alias(asset, assetAlias string) *BalanceStore {
	if asset != assetAlias {
		store.Aliases[asset] = assetAlias // alias
		store.Aliases[assetAlias] = asset // inverted alias
	}

	return store
}

func (store *BalanceStore) Aliased(asset string) (string, bool) {
	if alias, ok := store.Aliases[asset]; ok {
		return alias, true
	} else {
		return asset, false
	}
}
