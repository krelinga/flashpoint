package main

import (
	"cmp"
	"fmt"
	"maps"
	"math"
	"slices"

	"github.com/krelinga/flashpoint/go/dice"
)

type Simulation struct {
	Ships         []*Ship
	Distance      Distance
	Roller        dice.Roller
	Logger        Logger
	MaxRounds     int
	VolleyPlanner VolleyPlanner
}

type Result struct {
	Ships       []*Ship
	Rounds      int
	ParamValues map[string]string
}

func (s Simulation) RunOne() Result {
	if len(s.Ships) != 2 {
		panic("Simulation requires exactly 2 ships")
	}
	if s.Roller == nil {
		panic("Simulation requires a Roller")
	}
	if s.Distance == nil {
		panic("Simulation requires a non-nil distance")
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
		theRange := s.Distance.GetDistance(rounds)

		for i, ship := range turnOrder {
			if ship.IsDestroyed() {
				continue
			}
			roundLogger.Log(WeaponAttackTurnEvent{Ship: ship})
			turnLogger := roundLogger.Child()

			target := turnOrder[(i+1)%len(turnOrder)]
			volleys := volleyPlanner.PlanVolleys(ship, target, theRange)
			volleys.panicIfInvalidAiUsage(ship.AI)
			for _, volley := range volleys {
				s.fireVolley(ship, target, volley, theRange, turnLogger)
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

func (s Simulation) fireVolley(attacker, target *Ship, volley Volley, theRange int, logger Logger) {
	event := AttackEvent{
		Attacker: attacker,
		Target:   target,
		Range:    theRange,
		Volley:   volley,
	}
	weapon := volley.getWeapon(attacker)
	weaponkDice := weapon.Ranges.Find(theRange)
	if weaponkDice == nil {
		panic(fmt.Sprintf("weapon %s out of range %d", weapon.Name, theRange))
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

func (s Simulation) RunMany(n int, sink func(Result), params Params) {
	if params == nil {
		for range n {
			sink(s.RunOne())
		}
		return
	}

	params = maps.Clone(params)
	for key, values := range params {
		values = slices.Clone(values)
		slices.SortFunc(values, func(a, b Param) int {
			return cmp.Compare(a.String(), b.String())
		})
		params[key] = values
	}
	keys := make([]string, 0, len(params))
	keys = slices.AppendSeq(keys, maps.Keys(params))
	slices.Sort(keys)

	limits := make([]int, len(keys))
	for i := range limits {
		limits[i] = len(params[keys[i]])
	}
	var positions []int
	next := func() bool {
		if positions == nil {
			positions = make([]int, len(keys))
			return true
		}
		for i := len(positions) - 1; i >= 0; i-- {
			positions[i]++
			if positions[i] < limits[i] {
				return true
			}
			positions[i] = 0
		}
		return false
	}

	for next() {
		paramValues := map[string]string{}
		for i, position := range positions {
			key := keys[i]
			param := params[key][position]
			param.Apply(&s)
			paramValues[key] = param.String()
		}

		for range n {
			result := s.RunOne()
			result.ParamValues = paramValues
			sink(result)
		}
	}
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
