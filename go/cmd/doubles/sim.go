package main

type Simulation struct {
	Ships  []*Ship
	Roller Roller
	Logger Logger
}

func (s Simulation) Run() {
	if len(s.Ships) != 2 {
		panic("Simulation requires exactly 2 ships")
	}
	if s.Roller == nil {
		panic("Simulation requires a Roller")
	}
	if s.Logger == nil {
		s.Logger = NullLogger{}
	}

	done := func() bool {
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

	round := 0
	for !done() {
		round++
		s.Logger.Log(NewRoundEvent{Round: round})
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
			attackDice := Add(ship.LaserDice, ship.Crew)
			for range ship.NumLasers {
				event := AttackEvent{
					Attacker: ship,
					Target:   target,
				}
				if result := attackDice.Roll(s.Roller); result >= target.Aspect {
					damage := target.TakeDamage(ship.LaserDamage)
					event.Damage = &damage
				}
				turnLogger.Log(event)
			}
		}
	}
}
