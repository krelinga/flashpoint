package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/krelinga/flashpoint/go/dice"
	"github.com/krelinga/flashpoint/go/tables"
)

func main() {
	const janusName = "Janus"
	const griffinName = "Griffin"
	griffin := &Ship{
		Name:       griffinName,
		Crew:       dice.D8(),
		Aspect:     11,
		BarrierMax: 7,
		Armor:      3,
		Hull:       10,
		HullBoxes: map[int]Box{
			2: BoxStress,
			9: BoxCrit,
		},
		Core: 16,
		CoreBoxes: map[int]Box{
			1:  BoxCrit,
			3:  BoxCrit,
			8:  BoxStress,
			10: BoxCrit,
			14: BoxCrit,
		},
		OtherBarriersMax: []int{
			// front sides
			6, 6,
			// rear sides
			5, 5,
			// aft
			5,
		},
		StressMax:      6,
		StressRecovery: 2,
		Weapons: []*Weapon{
			WeaponLaserMk2,
			WeaponLaserMk2,
			WeaponLaserMk2,
			WeaponLaserMk2,
		},
		SalvosPerTurn:         1,
		SalvoAmmo:             5,
		SalvoThreat:           3,
		PointDefense:          dice.D8(),
		NumPointDefense:       2,
		PointDefenseThreshold: 6,
		AI:                    true,
		BarrierBoost:          true,
		Engineering:           MaxEngineering{},
	}

	janus := &Ship{
		Name:       janusName,
		Crew:       dice.D8(),
		Aspect:     11,
		BarrierMax: 10,
		Armor:      3,
		Hull:       15,
		HullBoxes: map[int]Box{
			4:  BoxStress,
			6:  BoxCrit,
			8:  BoxStress,
			12: BoxCrit,
		},
		Core: 12,
		CoreBoxes: map[int]Box{
			2: BoxCrit,
			4: BoxCrit,
			6: BoxCrit,
			9: BoxCrit,
		},
		OtherBarriersMax: []int{
			// front sides
			9, 9,
			// rear sides
			7, 7,
			// aft
			5,
		},
		StressMax:      7,
		StressRecovery: 3,
		Weapons: []*Weapon{
			WeaponLaserMk3,
			WeaponLaserMk3,
			WeaponLaserMk3,
			WeaponLaserMk3,
			WeaponPulseLaserMk2,
		},
		PointDefense:             dice.D8(),
		NumPointDefense:          3,
		PointDefenseThreshold:    6,
		BarrierBoost:             true,
		ImprovedBarrierGenerator: true,
		Engineering:              MaxEngineering{},
	}

	sim := Simulation{
		Ships:     []*Ship{griffin, janus},
		Roller:    dice.Random(rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))),
		MaxRounds: 10,
	}

	a := &Aggregator{
		Properties: map[string]func(Result) string{
			"Outcome": func(r Result) string {
				findShip := func(name string) *Ship {
					for _, ship := range r.Ships {
						if ship.Name == name {
							return ship
						}
					}
					panic(fmt.Sprintf("ship %s not found", name))
				}
				janus := findShip(janusName)
				griffin := findShip(griffinName)
				janusDestroyed := janus.IsDestroyed()
				griffinDestroyed := griffin.IsDestroyed()

				switch {
				case janusDestroyed && !griffinDestroyed:
					return "Griffin Wins"
				case !janusDestroyed && griffinDestroyed:
					return "Janus Wins"
				case janusDestroyed && griffinDestroyed:
					return "Both Destroyed"
				default:
					return "Both Survived"
				}
			},
		},
	}

	const iterations = 100_000
	sim.RunMany(iterations, a.Add, Params{
		"Distance": []Param{
			ConstantDistance(10),
			ClosingDistance{
				Start:       25,
				ClosingRate: 10,
				Min:         4,
			},
		},
		"VolleyPlanner": []Param{
			BasicVolleyPlanner{},
			PairingVolleyPlanner{},
		},
	})

	report := NewReport(a, ReportProperty("Outcome"), ReportParam("Distance"), func(key AggregateKey) bool {
		return key.HasParamValue("VolleyPlanner", "PairingVolleyPlanner")
	})
	reportTable := report.ToTable()
	fmt.Println(tables.PrettyString(reportTable, tables.PrettyOptions[string]{
		Title:        "Outcomes By Distance with Pairing volleys",
		UTF8:         true,
		Box:          true,
		HeaderRow:    true,
		HeaderColumn: true,
	}))
}
