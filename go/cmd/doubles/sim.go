package main

import (
	"fmt"
	"math"

	"github.com/krelinga/flashpoint/go/dice"
)

type Simulation struct {
	Ships     []*Ship
	Range     int
	Roller    dice.Roller
	Logger    Logger
	MaxRounds int
	VolleyPlanner VolleyPlanner
}

type Result struct {
	Ships  []*Ship
	Rounds int
}

func (s Simulation) RunOne() Result {
	if len(s.Ships) != 2 {
		panic("Simulation requires exactly 2 ships")
	}
	if s.Roller == nil {
		panic("Simulation requires a Roller")
	}
	if s.Range <= 0 {
		panic("Simulation requires a positive range")
	}

	var logger Logger = NullLogger{}
	if s.Logger != nil {
		logger = s.Logger
	}
	var volleyPlanner VolleyPlanner = BasicVolleyPlanner{}
	if s.VolleyPlanner != nil {
		volleyPlanner = s.VolleyPlanner
	}
	ships := make([]*Ship, 2)
	for i, ship := range s.Ships {
		ships[i] = ship.Clone()
	}

	salvoQueue := &SalvoQueue{}

	rounds := 0
	done := func() bool {
		if s.MaxRounds > 0 && rounds >= s.MaxRounds {
			return true
		}
		return ships[0].IsDestroyed() || ships[1].IsDestroyed()
	}

	rollInitiative := func() []*Ship {
		results := make([]int, len(ships))
		for results[0] == results[1] {
			for i, ship := range ships {
				results[i] = dice.Roll(ship.Crew, s.Roller).Total()
			}
		}
		if results[0] > results[1] {
			return ships
		} else {
			return []*Ship{ships[1], ships[0]}
		}
	}

	for !done() {
		rounds++
		logger.Log(NewRoundEvent{Round: rounds})
		roundLogger := logger.Child()
		turnOrder := rollInitiative()
		roundLogger.Log(InitiativeEvent{Winner: turnOrder[0]})

		salvosThisRound := salvoQueue.PopAll(rounds)
		for _, ship := range turnOrder {
			if ship.IsDestroyed() {
				continue
			}
			event := PDCTurnEvent{Target: ship}
			for _, salvo := range salvosThisRound {
				if salvo.Target == ship {
					event.BaseThreat += salvo.ThreatValue
				}
			}
			if event.BaseThreat == 0 {
				continue
			}
			for range ship.NumPointDefense {
				result := dice.Roll(ship.PointDefense, s.Roller).Total()
				if result >= ship.PointDefenseThreshold {
					event.PdcHits++
				}
			}
			if adjustedThreat := event.AdjustedThreat(); adjustedThreat > 0 {
				damage := SalvoDamage(adjustedThreat)
				event.Damage = ship.TakeDamage(damage)
			}
			roundLogger.Log(event)
		}

		for i, ship := range turnOrder {
			if ship.IsDestroyed() {
				continue
			}
			roundLogger.Log(WeaponAttackTurnEvent{Ship: ship})
			turnLogger := roundLogger.Child()

			target := turnOrder[(i+1)%len(turnOrder)]
			volleys := volleyPlanner.PlanVolleys(ship, target, s.Range)
			volleys.panicIfInvalidAiUsage(ship.AI)
			for _, volley := range volleys {
				s.fireVolley(ship, target, volley, turnLogger)
			}

			totalSalvos := 0
			for range ship.SalvosPerTurn {
				if ship.SalvoAmmo == 0 {
					continue
				}
				totalSalvos++
				salvoQueue.Push(&Salvo{
					Target:       target,
					ThreatValue:  ship.SalvoThreat,
					ArrivalRound: rounds + 1,
				})
				ship.SalvoAmmo--
			}
			if totalSalvos > 0 {
				turnLogger.Log(SalvoLaunchEvent{
					Attacker:    ship,
					Target:      target,
					SalvoThreat: ship.SalvoThreat,
					NumSalvos:   totalSalvos,
				})
			}
		}
	}
	return Result{
		Ships:  ships,
		Rounds: rounds,
	}
}

func (s Simulation) fireVolley(attacker, target *Ship, volley Volley, logger Logger) {
	event := AttackEvent{
		Attacker: attacker,
		Target: target,
		Range: s.Range,
		Volley: volley,
	}
	weapon := volley.getWeapon(attacker)
	weaponkDice := weapon.Ranges.Find(s.Range)
	if weaponkDice == nil {
		panic(fmt.Sprintf("weapon %s out of range %d", weapon.Name, s.Range))
	}
	attackDice := dice.Add(weaponkDice, attacker.Crew)
	if volley.UseAI {
		attackDice = dice.Add(attackDice, dice.D6())
	}
	results := dice.Roll(attackDice, s.Roller)
	event.AnyRerolls = handlePulseLasers(results, s.Roller)
	total := results.Total()
	event.Roll = total
	if total >= target.Aspect {
		event.Damage = target.TakeDamage(volley.Damage(attacker))
	}
	logger.Log(event)
}

func (s Simulation) RunMany(n int) []Result {
	results := make([]Result, n)
	for i := range n {
		results[i] = s.RunOne()
	}
	return results
}

func handlePulseLasers(results dice.ResultPool, roller dice.Roller) bool {
	var minIdx int
	minValue := math.MaxInt
	for idx, result := range results.Find(PulseTag) {
		if result.Value < minValue {
			minValue = result.Value
			minIdx = idx
		}
	}
	if minValue == math.MaxInt {
		return false
	}
	results[minIdx] = results[minIdx].Reroll(roller)
	return true
}