package main

type Distance interface {
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