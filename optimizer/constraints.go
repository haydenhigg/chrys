package optimizer

type Constraint = func(float64) float64

func Min(bound float64) Constraint {
	return func(x float64) float64 {
		if x < bound {
			return bound
		} else {
			return x
		}
	}
}

func Max(bound float64) Constraint {
	return func(x float64) float64 {
		if x > bound {
			return bound
		} else {
			return x
		}
	}
}

func applyConstraints(x float64, constraints []Constraint) float64 {
	for _, constraint := range constraints {
		x = constraint(x)
	}

	return x
}
