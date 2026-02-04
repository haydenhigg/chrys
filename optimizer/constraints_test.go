package optimizer

import (
	"testing"
)

func Test_Min(t *testing.T) {
	x := Min(0)(5)
	if x != 5 {
		t.Errorf("Min(0)(5) != 5: %f", x)
	}

	x = Min(0)(-1)
	if x != 0 {
		t.Errorf("Min(0)(-1) != 0: %f", x)
	}
}

func Test_Max(t *testing.T) {
	x := Max(10)(5)
	if x != 5 {
		t.Errorf("Max(10)(5) != 5: %f", x)
	}

	x = Max(10)(11)
	if x != 10 {
		t.Errorf("Min(10)(11) != 10: %f", x)
	}
}

func Test_applyConstraints(t *testing.T) {
	constraints := []Constraint{Min(0), Max(1)}

	x := applyConstraints(.1337, constraints)
	if x != .1337 {
		t.Errorf("applyConstraints(0.1337, constraints) != .1337: %f", x)
	}

	x = applyConstraints(-.1337, constraints)
	if x != 0 {
		t.Errorf("applyConstraints(-0.1337, constraints) != 0: %f", x)
	}

	x = applyConstraints(1.337, constraints)
	if x != 1 {
		t.Errorf("applyConstraints(1.337, constraints) != 1: %f", x)
	}
}
