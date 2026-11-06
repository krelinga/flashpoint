package main

type Ship struct {
	Crew Dice

	Aspect int
	Barrier int
	Armor int
	Hull int
	Core int

	NumLasers int
	LaserDamage int
	LaserDice Dice
}

func (s *Ship) TakeDamage(damage int) {
	if barrierDamage := min(damage, s.Barrier); barrierDamage > 0 {
		s.Barrier -= barrierDamage
		damage -= barrierDamage
	}

	if s.Armor >= damage {
		return
	} else {
		damage -= s.Armor
		s.Armor = max(0, s.Armor - 1)
	}

	if hullDamage := min(damage, s.Hull); hullDamage > 0 {
		s.Hull -= hullDamage
		damage -= hullDamage
	}

	if coreDamage := min(damage, s.Core); coreDamage > 0 {
		s.Core -= coreDamage
	}
}

func (s *Ship) IsDestroyed() bool {
	return s.Core <= 0
}
