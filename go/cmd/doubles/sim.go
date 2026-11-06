package main

type Simulation struct {
	Ships  []*Ship
	Roller Roller
	Logger Logger
	MaxRounds int
}

func (s Simulation) Run() int {
	if len(s.Ships) != 2 {
		panic("Simulation requires exactly 2 ships")
	}
	if s.Roller == nil {
		panic("Simulation requires a Roller")
	}
	if s.Logger == nil {
		s.Logger = NullLogger{}
	}

	rounds := 0
	done := func() bool {
		if s.MaxRounds > 0 && rounds >= s.MaxRounds {
			return true
		}
		return s.Ships[0].IsDestroyed() || s.Ships[1].IsDestroyed()
	}

	rollInitiative := func() []*Ship {
		results := make([]int, len(s.Ships))
		for results[0] == results[1] {
			for i, ship := range s.Ships {
				results[i] = ship.Crew.Roll(s.Roller)
			}
		}
		if results[0] > results[1] {
			return s.Ships
		} else {
			return []*Ship{s.Ships[1], s.Ships[0]}
		}
	}

	for !done() {
		rounds++
		s.Logger.Log(NewRoundEvent{Round: rounds})
		roundLogger := s.Logger.Child()
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
					AI : ai,
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
	return rounds
}
