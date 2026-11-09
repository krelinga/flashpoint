package main

import (
	"cmp"
	"slices"
)

type Range struct {
	Limit int
	Dice Dice
}

type Ranges []Range

func (r Ranges) Find(wantRange int) Dice {
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