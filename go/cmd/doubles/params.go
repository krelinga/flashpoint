package main

import "fmt"

type Param interface {
	Apply(*Simulation)
	String() string
}

type Params map[string][]Param

func (p Params) PanicIfInvalid() {
	if len(p) == 0 {
		panic("Params cannot be empty")
	}
	for key, vals := range p {
		if len(vals) == 0 {
			panic(fmt.Sprintf("no values for Params key %s", key))
		}
	}
}

func (p Params) NumVariations() int {
	p.PanicIfInvalid()
	var variations int
	first := true
	for _, vals := range p {
		if first {
			variations = len(vals)
			first = false
		} else {
			variations *= len(vals)
		}
	}
	return variations
}