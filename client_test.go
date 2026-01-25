package chrys

import (
	"github.com/haydenhigg/chrys/frame"
	"math"
	"testing"
	"time"
)

// mock
type MockAPI struct {
	callback func()
}

func (api MockAPI) FetchBalances() (map[string]float64, error) {
	return map[string]float64{
		"USD": 133.7,
		"BTC": 0.001337,
		"ETH": 0.01337,
	}, nil
}

func (api MockAPI) FetchFramesSince(
	pair string,
	interval time.Duration,
	since time.Time,
) ([]*frame.Frame, error) {
	frames := []*frame.Frame{}
	i := 0.

	price := 0.
	switch pair {
	case "BTC/USD":
		price = 88304.55
	case "ETH/USD":
		price = 2943.89
	}

	// all frames that start !Before(since) and end !After(now-interval)
	now := time.Now().Add(-interval)
	for t := since; !t.After(now); t = t.Add(interval) {
		if t.Equal(t.Truncate(interval)) {
			frames = append(frames, &frame.Frame{
				Time:  t,
				Close: price + i,
			})
			i++
		} else {
			t = t.Truncate(interval)
		}
	}

	return frames, nil
}

func (api MockAPI) MarketOrder(side, pair string, quantity float64) error {
	return nil
}

// tests
func Test_SetFee(t *testing.T) {
	// create Client
	client := NewClient(nil)

	// SetFee()
	client.SetFee(0.01337)

	// assert
	if client.Fee != 0.01337 {
		t.Errorf("client.Fee != 0.01337: %f", client.Fee)
	}
}

func Test_SetIsLive(t *testing.T) {
	// create Client
	client := NewClient(nil)

	// assert default
	if client.IsLive {
		t.Errorf("client.Fee != false: %v", client.IsLive)
	}

	// SetIsLive()
	client.SetIsLive(true)

	// assert
	if !client.IsLive {
		t.Errorf("client.Fee != true: %v", client.IsLive)
	}
}

func Test_Value(t *testing.T) {
	// create Client
	client := NewClient(MockAPI{})

	// Value()
	value, err := client.Value([]string{"USD", "ETH", "BTC"}, time.Now())
	if err != nil {
		t.Errorf("err: %v", err)
	}

	// assert
	if math.Abs(value-291.12299265) > 10e-6 {
		t.Errorf("value != 291.12299265: %f", value)
	}

}

func Test_ValueAliases(t *testing.T) {
	// create Client
	client := NewClient(MockAPI{})

	// set aliases
	client.Balances.
		Alias("BTC", "XXBT").
		Alias("ZETH", "ETH").
		Alias("ZUSD", "USD")

	// Value()
	value, err := client.Value([]string{"USD", "ETH", "BTC"}, time.Now())
	if err != nil {
		t.Errorf("err: %v", err)
	}

	// assert
	if math.Abs(value-291.12299265) > 10e-6 {
		t.Errorf("value != 291.12299265: %f", value)
	}
}
