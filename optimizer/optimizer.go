package optimizer

import "maps"

type Input = map[string]float64
type Constraints = map[string][]Constraint
type ObjectiveFunc = func(Input) float64

type Optimizer struct {
	F           ObjectiveFunc
	x           Input
	constraints Constraints
}

// initializer
func New(objective ObjectiveFunc) *Optimizer {
	return &Optimizer{
		F:           objective,
		x:           Input{},
		constraints: Constraints{},
	}
}

// constraints
func (opt *Optimizer) Constrain(
	k string,
	newConstraints ...Constraint,
) *Optimizer {
	if constraints, ok := opt.constraints[k]; !ok {
		opt.constraints[k] = newConstraints
	} else {
		opt.constraints[k] = append(constraints, newConstraints...)
	}

	return opt
}

func (opt *Optimizer) withConstraints(k string, v float64) float64 {
	if constraints, ok := opt.constraints[k]; ok {
		return applyConstraints(v, constraints)
	}

	return v
}

// X setter
func (opt *Optimizer) SetX(x Input) *Optimizer {
	for k, v := range x {
		opt.x[k] = opt.withConstraints(k, v)
	}

	return opt
}

// X getter
func (opt *Optimizer) X() Input {
	x := make(Input, len(opt.x))
	maps.Copy(x, opt.x)

	return x
}
