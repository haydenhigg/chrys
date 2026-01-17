package algo

import domain "github.com/haydenhigg/chrys/frame"

type Machine interface {
	Apply(x float64) Machine
	ApplyFrame(frame *domain.Frame) Machine
	Val() float64
}

type Composer struct {
	Machines []Machine
	Value    float64
}

func NewComposer(initial Machine) *Composer {
	return &Composer{
		Machines: []Machine{initial},
	}
}

func (composer *Composer) Of(machine Machine) *Composer {
	composer.Machines = append(composer.Machines, machine)
	return composer
}

func (composer *Composer) feedForward(initial func(machine Machine) float64) {
	endIndex := len(composer.Machines) - 1
	for i := range composer.Machines {
		machine := composer.Machines[endIndex-i]
		if i == 0 {
			composer.Value = initial(machine)
		} else {
			composer.Value = machine.Apply(composer.Value).Val()
		}
	}
}

func (composer *Composer) Apply(x float64) Machine {
	composer.feedForward(func(machine Machine) float64 {
		return machine.Apply(x).Val()
	})

	return composer
}

func (composer *Composer) ApplyFrame(frame *domain.Frame) Machine {
	composer.feedForward(func(machine Machine) float64 {
		return machine.ApplyFrame(frame).Val()
	})

	return composer
}

func (composer *Composer) Val() float64 {
	return composer.Value
}
