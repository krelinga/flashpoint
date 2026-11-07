package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	const janusName = "Janus"
	const griffinName = "Griffin"
	griffin := &Ship{
		Name:                  griffinName,
		Crew:                  Sides(8),
		Aspect:                11,
		Barrier:               7,
		Armor:                 3,
		Hull:                  10,
		Core:                  16,
		NumLasers:             2,
		LaserDamage:           6,
		LaserDice:             Sides(8),
		SalvosPerTurn:         1,
		SalvoAmmo:             5,
		SalvoThreat:           3,
		PointDefense:          Sides(8),
		NumPointDefense:       2,
		PointDefenseThreshold: 6,
		AI:                    true,
	}

	janus := &Ship{
		Name:                  janusName,
		Crew:                  Sides(8),
		Aspect:                11,
		Barrier:               10,
		Armor:                 3,
		Hull:                  15,
		Core:                  12,
		NumLasers:             2,
		LaserDamage:           8,
		PointDefense:          Sides(8),
		NumPointDefense:       3,
		PointDefenseThreshold: 6,
		LaserDice:             Sides(10),
	}

	sim := Simulation{
		Ships:     []*Ship{griffin, janus},
		Roller:    &RandomRoller{Rand: rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))},
		MaxRounds: 10,
	}

	const iterations = 100_000
	results := sim.RunMany(iterations)

	findShip := func(ships []*Ship, name string) *Ship {
		for _, s := range ships {
			if s.Name == name {
				return s
			}
		}
		panic(fmt.Sprintf("No ship named %s", name))
	}

	const janusWins = "Janus Wins"
	const griffinWins = "Griffin Wins"
	const bothDestroyed = "Both Destroyed"
	const bothSurvive = "Both Survive"
	widenedResults := Widen(results,
		Predicate{
			Name: janusWins,
			Test: func(r Result) bool {
				janus := findShip(r.Ships, janusName)
				griffin := findShip(r.Ships, griffinName)
				return !janus.IsDestroyed() && griffin.IsDestroyed()
			},
		},
		Predicate{
			Name: griffinWins,
			Test: func(r Result) bool {
				janus := findShip(r.Ships, janusName)
				griffin := findShip(r.Ships, griffinName)
				return !griffin.IsDestroyed() && janus.IsDestroyed()
			},
		},
		Predicate{
			Name: bothDestroyed,
			Test: func(r Result) bool {
				for _, ship := range r.Ships {
					if !ship.IsDestroyed() {
						return false
					}
				}
				return true
			},
		},
		Predicate{
			Name: bothSurvive,
			Test: func(r Result) bool {
				for _, ship := range r.Ships {
					if ship.IsDestroyed() {
						return false
					}
				}
				return true
			},
		},
	)
	matchers := []Matcher{
		Contains(janusWins),
		Contains(griffinWins),
		Contains(bothDestroyed),
		Contains(bothSurvive),
	}
	for _, m := range matchers {
		fmt.Println(m, widenedResults.Compute(m))
	}
}
