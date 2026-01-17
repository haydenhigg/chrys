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
	mockAPI := MockBalanceAPI{
		callback: func() {
			didUseAPI = true
		},
	}

	// Get()
	balances, err := NewBalances(mockAPI).Get()
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	assertBalancesEqual(map[string]float64{
		"USD": 133.7,
		"BTC": 0.001337,
		"ETH": 0.01337,
	}, balances, t)

	if !didUseAPI {
		t.Errorf("cache was hit")
	}
}

func Test_GetCached(t *testing.T) {
	// set up mock
	didUseAPI := false
	mockAPI := MockBalanceAPI{
		callback: func() {
			didUseAPI = true
		},
	}

	// Get()
	store := NewBalances(mockAPI)
	store.Get()

	// reset didUseAPI
	didUseAPI = false

	// Get() again
	balances, err := store.Get()
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}

	// assert
	assertBalancesEqual(map[string]float64{
		"USD": 133.7,
		"BTC": 0.001337,
		"ETH": 0.01337,
	}, balances, t)

	if didUseAPI {
		t.Errorf("cache was not hit")
	}
}

func Test_Set(t *testing.T) {
	// set up mock
	mockAPI := MockBalanceAPI{}

	store := NewBalances(mockAPI)

	// Set()
	targetBalances, _ := mockAPI.FetchBalances()
	store.Set(targetBalances)

	// assert
	assertBalancesEqual(map[string]float64{
		"USD": 133.70,
		"BTC": 0.001337,
		"ETH": 0.01337,
	}, store.balances, t)
}

func Test_SetAddSubtract(t *testing.T) {
	// set up mock
	mockAPI := MockBalanceAPI{}

	store := NewBalances(mockAPI)

	// Set()
	targetBalances, _ := mockAPI.FetchBalances()
	store.Set(targetBalances)

	// Set() with existing keys
	store.Set(map[string]float64{
		"USD": -43.94,
		"ETH": 0.01337,
	})

	// assert
	assertBalancesEqual(map[string]float64{
		"USD": 89.76,
		"BTC": 0.001337,
		"ETH": 0.02674,
	}, store.balances, t)
}

func Test_AliasSet(t *testing.T) {
	// set up mock
	mockAPI := MockBalanceAPI{}

	store := NewBalances(mockAPI)
	store.Alias("BTC", "XBT.F") // alias
	store.Alias("ZUSD", "USD")  // inverted alias

	// Set()
	targetBalances, _ := mockAPI.FetchBalances()
	store.Set(targetBalances)

	// assert
	assertBalancesEqual(map[string]float64{
		"USD":   133.70,
		"ZUSD":  133.70,
		"BTC":   0.001337,
		"XBT.F": 0.001337,
		"ETH":   0.01337,
	}, store.balances, t)
}

func Test_AliasSetAddSubtract(t *testing.T) {
	// set up mock
	mockAPI := MockBalanceAPI{}

	store := NewBalances(mockAPI)
	store.Alias("BTC", "XBT.F") // normal alias
	store.Alias("ZUSD", "USD")  // inverted alias

	// Set()
	targetBalances, _ := mockAPI.FetchBalances()
	store.Set(targetBalances)

	// Set() with existing keys
	store.Set(map[string]float64{
		"USD": 32.13,
		"BTC": -0.000337,
	})

	// assert
	assertBalancesEqual(map[string]float64{
		"USD":   165.83,
		"ZUSD":  165.83,
		"BTC":   0.001,
		"XBT.F": 0.001,
		"ETH":   0.01337,
	}, store.balances, t)
}
