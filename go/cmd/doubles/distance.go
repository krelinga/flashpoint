package main

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

func (cd *ClosingDistance) GetDistance(round int) int {
	if cd.Start <= 0 {
		panic("Start must be positive")
	}
	if cd.Min < 0 {
		panic("Min must be non-negative")
	}
	theMin := max(cd.Min, 1)

	return max(theMin, cd.Start-cd.ClosingRate*(round-1))
}
