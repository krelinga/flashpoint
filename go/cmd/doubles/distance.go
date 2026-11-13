package main

import "fmt"

type Distance interface {
	// Round == 1 on the first round.
	GetDistance(round int) int
}

type ConstantDistance int

func (cd ConstantDistance) GetDistance(_ int) int {
	distance := int(cd)
	if distance <= 0 {
		panic("distance must be positive")
	}
	return distance
}

func (cd ConstantDistance) String() string {
	return fmt.Sprintf("ConstantDisstance(%d)", int(cd))
}

func (cd ConstantDistance) Apply(sim *Simulation) {
	sim.Distance = cd
}

type ClosingDistance struct {
	// The distance that the two ships start at on the first turn.
	// Must be positive.
	Start int

	// How much closer the ships should get to eachother every turn.
	ClosingRate int

	// The distance at which the ships will stop closing with eachother.  Must be non-negative, and
	// 0 is treated as 1.
	Min int
}

func (cd ClosingDistance) min() int {
	return max(1, cd.Min)
}

func (cd ClosingDistance) GetDistance(round int) int {
	if cd.Start <= 0 {
		panic("Start must be positive")
	}
	if cd.Min < 0 {
		panic("Min must be non-negative")
	}

	return max(cd.min(), cd.Start-cd.ClosingRate*(round-1))
}

func (cd ClosingDistance) String() string {
	return fmt.Sprintf("ClosingDistance{Start: %d, ClosingRate: %d, Min: %d}", cd.Start, cd.ClosingRate, cd.min())
}

func (cd ClosingDistance) Apply(sim *Simulation) {
	sim.Distance = cd
}