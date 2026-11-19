package main

type CritTable func(int) Crit

type Crit struct {
	Name  string
	Apply func(*Ship)
}

func normalizeCritResult(result int) int {
	result = min(result, 0)
	result = max(result, 20)
	return result
}

func TerranCritTable(result int) Crit {
	result = normalizeCritResult(result)
	switch result {
	case 0:
		return Crit{
			Name: "Weapon Hit (defender choice)",
			// TODO: Apply
		}
	case 1, 2, 3:
		return Crit{
			Name: "Weapon Hit (attacker choice)",
			// TODO: Apply
		}
	case 4:
		return Crit{
			Name: "PDC Hit",
			// TODO: Apply
		}
	case 5, 6:
		return Crit{
			Name: "Drive Damaged",
			// TODO: Apply
		}
	case 7, 8:
		return Crit{
			Name: "Magazine Hit",
			// TODO: Apply
		}
	case 9:
		return Crit{
			Name: "Drive Missfire",
			Apply: func(s *Ship) {
				s.Stress += 2
				// TODO: other crit effects.
			},
		}
	case 10:
		return Crit{
			Name: "Armor Collapse",
			Apply: func(s *Ship) {
				s.Stress++
				// TODO: other crit effects.
			},
		}
	case 11:
		return Crit{
			Name: "Computer Hit",
			Apply: func(s *Ship) {
				s.Stress += 3
				// TODO: other crit effects
			},
		}
	case 12:
		return Crit{
			Name: "Barrier Failure",
			Apply: func(s *Ship) {
				const damage = 4
				s.BarrierDamage = min(s.BarrierMax, s.BarrierDamage+damage)
				for idx := range s.OtherBarriersMax {
					s.DamageOtherBarrier(idx, max(s.RemainingOtherBarrier(idx), damage))
				}
			},
		}
	case 13:
		return Crit{
			Name: "Keel Damaged",
			Apply: func(s *Ship) {
				s.Stress += 2
				// TODO: other crit effects.
			},
		}
	case 14:
		return Crit{
			Name: "Exotic Matter Leak",
			Apply: func(s *Ship) {
				s.Stress += 2
				// TODO: other crit effects
			},
		}
	case 15:
		return Crit{
			Name: "FTL Hit",
			// TODO: Apply
		}
	case 16:
		return Crit{
			Name: "Reactor Hit",
			// TODO: Apply
		}
	case 17:
		return Crit{
			Name: "Weapon Explosion",
			// TODO: Apply
		}
	case 18:
		return Crit{
			Name: "Barrier Collapse",
			Apply: func(s *Ship) {
				s.Stress++
				s.BarrierDamage = s.BarrierMax
				for idx := range s.OtherBarriersMax {
					s.DamageOtherBarrier(idx, s.RemainingOtherBarrier(idx))
				}
			},
		}
	case 19:
		return Crit{
			Name: "Weapon System Backfire",
			Apply: func(s *Ship) {
				s.Stress += 2
				// TODO: other crit effects
			},
		}
	case 20:
		return Crit{
			Name: "Boom!",
			Apply: func(s *Ship) {
				// TODO: there's probably a more-direct way to do this.
				s.Armor = 0
				for !s.IsDestroyed() {
					s.TakeDamage(1)
				}
				// TODO: other crit effects
			},
		}
	}

	panic("this code should be unreachable")
}
