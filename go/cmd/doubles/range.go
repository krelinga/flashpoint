package main

import (
	"cmp"
	"slices"

	"github.com/krelinga/flashpoint/go/dice"
)

type Range struct {
	Limit int
	Dice dice.ModelPool
}

type Ranges []Range

func (r Ranges) Find(wantRange int) dice.ModelPool {
	if !slices.IsSortedFunc(r, func(a, b Range) int {
		return cmp.Compare(a.Limit, b.Limit)
	}) {
		panic("Ranges must be sorted by Limit")
	}

	for _, entry := range r {
		if wantRange <= entry.Limit {
			return entry.Dice
		}
	}

	return nil
}