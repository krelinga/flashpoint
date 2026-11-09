package dice

import "math/rand/v2"

type Roller interface {
	Roll(sides int) int
}

func Random(r *rand.Rand) Roller {
	if r == nil {
		panic("r must not be nil")
	}
	return randomRoller{r}
}

type randomRoller struct {
	r *rand.Rand
}

func (r randomRoller) Roll(sides int) int {
	return r.r.IntN(sides) + 1
}