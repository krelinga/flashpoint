package main

import "math/rand/v2"

func main() {
	griffin := &Ship{
		Crew: Multiple(2, Sides(6)),
		Aspect: 11,
		Barrier: 7,
		Armor: 3,
		Hull: 10,
		Core: 16,
		NumLasers: 4,
		LaserDamage: 3,
		LaserDice: Sides(8),
	}

	janus := &Ship{
		Crew: Sides(2),
		Aspect: 11,
		Barrier: 10,
		Armor: 3,
		Hull: 15,
		Core: 12,
		NumLasers: 4,
		LaserDamage: 4,
		LaserDice: Sides(2),
	}

	sim := &Simulation{
		Ships: []*Ship{griffin, janus},
		Roller: &RandomRoller{Rand: rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))},
	}

	sim.Run()

	if griffin.IsDestroyed() && !janus.IsDestroyed() {
		println("Janus wins!")
	} else if janus.IsDestroyed() && !griffin.IsDestroyed() {
		println("Griffin wins!")
	} else {
		println("Both ships destroyed!")
	}
}