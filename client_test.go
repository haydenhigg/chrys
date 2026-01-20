package chrys

import (
	// "github.com/haydenhigg/chrys/driver"
	"testing"
)

// mock
type MockOrdererAPI struct {
	callback func()
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
