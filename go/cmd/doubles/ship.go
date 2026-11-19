package main

import (
	"fmt"
	"strings"

	"github.com/krelinga/flashpoint/go/dice"
)

type Box bool

const (
	BoxStress Box = true
	BoxCrit   Box = false
)

type Ship struct {
	Name string
	Crew dice.ModelPool

	Aspect              int
	BarrierMax          int
	BarrierDamage       int
	Armor               int
	Hull                int
	HullBoxes           map[int]Box
	Core                int
	CoreBoxes           map[int]Box
	OtherBarriersMax    []int
	OtherBarriersDamage []int

	Stress         int
	StressMax      int
	StressRecovery int

	Weapons []*Weapon

	SalvosPerTurn int
	SalvoAmmo     int
	SalvoThreat   int

	PointDefense          dice.ModelPool
	NumPointDefense       int
	PointDefenseThreshold int

	// Abilities
	AI           bool
	BarrierBoost bool

	// Traits
	ImprovedBarrierGenerator bool

	Engineering Engineering
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
	Stress         int
	Crits          int
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
	if d.Stress > 0 {
		parts = append(parts, fmt.Sprintf("Stress: %d", d.Stress))
	}
	if d.Crits > 0 {
		parts = append(parts, fmt.Sprintf("Crits: %d", d.Crits))
	}
	return strings.Join(parts, " ")
}

func (s *Ship) checkOtherBarrierIdx(idx int) {
	if idx < 0 || idx >= len(s.OtherBarriersMax) {
		panic("idx out of range")
	}
}

func (s *Ship) GetOtherBarrierDamage(idx int) int {
	s.checkOtherBarrierIdx(idx)
	if s.OtherBarriersDamage == nil {
		return 0
	}
	return s.OtherBarriersDamage[idx]
}

func (s *Ship) RemainingOtherBarrier(idx int) int {
	s.checkOtherBarrierIdx(idx)
	return s.OtherBarriersMax[idx] - s.GetOtherBarrierDamage(idx)
}

func (s *Ship) DamageOtherBarrier(idx int, damage int) {
	s.checkOtherBarrierIdx(idx)
	if damage < 0 || damage > s.RemainingOtherBarrier(idx) {
		panic("damage out of range")
	}
	if s.OtherBarriersDamage == nil {
		s.OtherBarriersDamage = make([]int, len(s.OtherBarriersMax))
	}
	s.OtherBarriersDamage[idx] += damage
}

func (s *Ship) RemainingBarrier() int {
	return s.BarrierMax - s.BarrierDamage
}

func (s *Ship) RemainingStress() int {
	return s.StressMax - s.Stress
}

func (s *Ship) TakeDamage(damage int) Damage {
	d := Damage{
		Destroyed: s.IsDestroyed(),
	}
	if newBarrierDamage := min(damage, s.RemainingBarrier()); newBarrierDamage > 0 {
		s.BarrierDamage += newBarrierDamage
		damage -= newBarrierDamage
		d.Barrier = newBarrierDamage
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
