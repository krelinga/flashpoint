package main

import (
	"fmt"
	"strings"
)

type Ship struct {
	Name string
	Crew Dice

	Aspect  int
	Barrier int
	Armor   int
	Hull    int
	Core    int

	NumLasers   int
	LaserDamage int
	LaserDice   Dice

	AI bool
}

func (s *Ship) Clone() *Ship {
	clone := *s
	return &clone
}

type Damage struct {
	Barrier        int
	Armor          int
	ArmorReduction int
	Hull           int
	Core           int
	Destroyed      bool
}

func (d Damage) String() string {
	if d.Destroyed {
		return "Destroyed"
	}
	parts := []string{}
	if d.Barrier > 0 {
		parts = append(parts, fmt.Sprintf("Barrier: %d", d.Barrier))
	}
	if d.Armor > 0 {
		parts = append(parts, fmt.Sprintf("Armor: %d", d.Armor))
		if d.ArmorReduction == 0 {
			parts = append(parts, "(Bounces off armor)")
		} else {
			parts = append(parts, fmt.Sprintf("(Armor reduced by %d)", d.ArmorReduction))
		}
	}
	if d.Hull > 0 {
		parts = append(parts, fmt.Sprintf("Hull: %d", d.Hull))
	}
	if d.Core > 0 {
		parts = append(parts, fmt.Sprintf("Core: %d", d.Core))
	}
	return strings.Join(parts, " ")
}

func (s *Ship) TakeDamage(damage int) Damage {
	d := Damage{
		Destroyed: s.IsDestroyed(),
	}
	if barrierDamage := min(damage, s.Barrier); barrierDamage > 0 {
		s.Barrier -= barrierDamage
		damage -= barrierDamage
		d.Barrier = barrierDamage
	}

	if s.Armor >= damage {
		d.Armor = damage
		return d
	} else {
		damage -= s.Armor
		d.Armor = s.Armor
		oldArmor := s.Armor
		s.Armor = max(0, s.Armor-1)
		d.ArmorReduction = oldArmor - s.Armor
	}

	if hullDamage := min(damage, s.Hull); hullDamage > 0 {
		s.Hull -= hullDamage
		damage -= hullDamage
		d.Hull = hullDamage
	}

	if coreDamage := min(damage, s.Core); coreDamage > 0 {
		s.Core -= coreDamage
		d.Core = coreDamage
	}

	d.Destroyed = s.IsDestroyed()
	return d
}

func (s *Ship) IsDestroyed() bool {
	return s.Core <= 0
}
