package main

type Simulation struct {
	Ships     []*Ship
	Roller    Roller
	Logger    Logger
	MaxRounds int
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

	var logger Logger = NullLogger{}
	if s.Logger != nil {
		logger = s.Logger
	}
	ships := make([]*Ship, 2)
	for i, ship := range s.Ships {
		ships[i] = ship.Clone()
	}

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
				results[i] = ship.Crew.Roll(s.Roller)
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
		for i, ship := range turnOrder {
			if ship.IsDestroyed() {
				continue
			}
			roundLogger.Log(NewTurnEvent{Ship: ship})
			turnLogger := roundLogger.Child()

			target := turnOrder[(i+1)%len(turnOrder)]
			for n := range ship.NumLasers {
				ai := ship.AI && n == 0
				event := AttackEvent{
					Attacker: ship,
					Target:   target,
					AI:       ai,
				}
				attackDice := Add(ship.LaserDice, ship.Crew)
				if ai {
					attackDice = Add(attackDice, Sides(6))
				}
				result := attackDice.Roll(s.Roller)
				event.Roll = result
				if result >= target.Aspect {
					event.Damage = target.TakeDamage(ship.LaserDamage)
				}
				turnLogger.Log(event)
			}
		}
	}
	return Result{
		Ships:  ships,
		Rounds: rounds,
	}
}

func (s Simulation) RunMany(n int) []Result {
	results := make([]Result, n)
	for i := range n {
		results[i] = s.RunOne()
	}
	return results
}
