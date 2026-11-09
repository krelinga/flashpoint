package main

type Simulation struct {
	Ships     []*Ship
	Range     int
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
	if s.Range <= 0 {
		panic("Simulation requires a positive range")
	}

	var logger Logger = NullLogger{}
	if s.Logger != nil {
		logger = s.Logger
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
				result := ship.PointDefense.Roll(s.Roller)
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
			var usedAi bool
			for _, weapon := range ship.Weapons {
				laserDice := weapon.Ranges.Find(s.Range)
				if laserDice == nil {
					// out of range
					continue
				}
				ai := ship.AI && !usedAi
				if ai {
					usedAi = true
				}
				event := AttackEvent{
					Attacker: ship,
					Target:   target,
					Range:    s.Range,
					AI:       ai,
				}
				attackDice := Add(laserDice, ship.Crew)
				if ai {
					attackDice = Add(attackDice, Sides(6))
				}
				result := attackDice.Roll(s.Roller)
				event.Roll = result
				if result >= target.Aspect {
					event.Damage = target.TakeDamage(weapon.Damage)
				}
				turnLogger.Log(event)
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

func (s Simulation) RunMany(n int) []Result {
	results := make([]Result, n)
	for i := range n {
		results[i] = s.RunOne()
	}
	return results
}
