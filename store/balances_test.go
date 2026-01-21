package store

import (
	"math"
	"testing"
)

// mock
type MockBalanceAPI struct {
	callback func()
}

func (api MockBalanceAPI) FetchBalances() (map[string]float64, error) {
	if api.callback != nil {
		api.callback()
	}

	return map[string]float64{
		"USD": 133.7,
		"BTC": 0.001337,
		"ETH": 0.01337,
	}, nil
}

// tests
func assertBalancesEqual(a, b map[string]float64, t *testing.T) {
	for k, va := range a {
		if vb, ok := b[k]; !ok {
			t.Errorf(`b["%s"] does not exist`, k)
		} else if math.Abs(va-vb) > 10e-6 {
			t.Errorf(`a["%s"] != b["%s"]: %v != %v`, k, k, va, vb)
		}
	}

	for k := range b {
		if _, ok := a[k]; !ok {
			t.Errorf(`a["%s"] does not exist`, k)
		}
	}
}

func Test_GetUncached(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockBalanceAPI{callback: func() { didUseAPI = true }}

	// Get()
	balances, err := NewBalances(mockAPI).Get()
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	assertBalancesEqual(balances, map[string]float64{
		"USD": 133.7,
		"BTC": 0.001337,
		"ETH": 0.01337,
	}, t)

	if !didUseAPI {
		t.Errorf("cache was hit")
	}
}

func Test_GetCached(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockBalanceAPI{callback: func() { didUseAPI = true }}

	// set up store
	store := NewBalances(mockAPI)

	// Get()
	store.Get()

	// reset didUseAPI
	didUseAPI = false

	// Get() again
	balances, err := store.Get()
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	assertBalancesEqual(balances, map[string]float64{
		"USD": 133.7,
		"BTC": 0.001337,
		"ETH": 0.01337,
	}, t)

	if didUseAPI {
		t.Errorf("cache was not hit")
	}
}

func Test_Set_balances(t *testing.T) {
	// set up store
	store := NewBalances(MockBalanceAPI{})

	// Set()
	expectedBalances := map[string]float64{
		"USD": 133.70,
		"BTC": 0.001337,
		"ETH": 0.01337,
	}
	store.Set(expectedBalances)

	// assert
	assertBalancesEqual(store.balances, expectedBalances, t)
}

func Test_SetAddSubtract(t *testing.T) {
	// set up store
	store := NewBalances(MockBalanceAPI{})

	// Set()
	store.Set(map[string]float64{
		"USD": 133.70,
		"BTC": 0.001337,
		"ETH": 0.01337,
	})

	// Set() with existing keys
	store.Set(map[string]float64{
		"USD": -43.94,
		"ETH": 0.01337,
	})

	// assert
	assertBalancesEqual(store.balances, map[string]float64{
		"USD": 89.76,
		"BTC": 0.001337,
		"ETH": 0.02674,
	}, t)
}

func Test_AliasSet(t *testing.T) {
	// set up store
	store := NewBalances(MockBalanceAPI{})
	store.Alias("BTC", "XXBT") // alias
	store.Alias("ZUSD", "USD")  // inverted alias

	// Set()
	store.Set(map[string]float64{
		"USD": 133.70,
		"BTC": 0.001337,
		"ETH": 0.01337,
	})

	// assert
	assertBalancesEqual(store.balances, map[string]float64{
		"USD":   133.70,
		"ZUSD":  133.70,
		"BTC":   0.001337,
		"XXBT": 0.001337,
		"ETH":   0.01337,
	}, t)
}

func Test_AliasSetAddSubtract(t *testing.T) {
	// set up store
	store := NewBalances(MockBalanceAPI{})
	store.Alias("BTC", "XXBT") // normal alias
	store.Alias("ZUSD", "USD")  // inverted alias

	// Set()
	store.Set(map[string]float64{
		"USD": 133.70,
		"BTC": 0.001337,
		"ETH": 0.01337,
	})

	// Set() with existing keys
	store.Set(map[string]float64{
		"USD": 32.13,
		"BTC": -0.000337,
	})

	// assert
	assertBalancesEqual(store.balances, map[string]float64{
		"USD":   165.83,
		"ZUSD":  165.83,
		"BTC":   0.001,
		"XXBT": 0.001,
		"ETH":   0.01337,
	}, t)
}

func Test_Aliased(t *testing.T) {
	// set up store
	store := NewBalances(MockBalanceAPI{})
	store.Alias("BTC", "XXBT") // alias
	store.Alias("ZUSD", "USD")  // inverted alias

	// Aliased()
	alias, ok := store.Aliased("XXBT")
	if !ok {
		t.Errorf(`"ZUSD" is not aliased`)
	} else if alias != "BTC" {
		t.Errorf(`alias != "BTC": %s`, alias)
	}

	// inverted Aliased()
	alias, ok = store.Aliased("ZUSD")
	if !ok {
		t.Errorf(`"ZUSD" is not aliased`)
	} else if alias != "USD" {
		t.Errorf(`alias != "USD": %s`, alias)
	}
}
