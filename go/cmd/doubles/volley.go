package main

type Volley struct {
	WeaponIdxs []int
	UseAI      bool
}

func (v Volley) getWeapon(ship *Ship) *Weapon {
	if len(v.WeaponIdxs) == 0 {
		panic("volley.WeaponIdxs must not be empty")
	}
	if len(ship.Weapons) == 0 {
		panic("ship has no weapons")
	}
	weapon := ship.Weapons[v.WeaponIdxs[0]]
	for i := 1; i < len(v.WeaponIdxs); i++ {
		weaponIdx := v.WeaponIdxs[i]
		if weapon != ship.Weapons[weaponIdx] {
			panic("a Volley can only contain a single type of weapon")
		}
	}
	return weapon
}

func (v Volley) Damage(ship *Ship) int {
	weapon := v.getWeapon(ship)
	return weapon.Damage * len(v.WeaponIdxs)
}

type Volleys []Volley

func (v Volleys) panicIfInvalidAiUsage(shipHasAi bool) {
	aiAvail := shipHasAi
	for _, volley := range v {
		if volley.UseAI && !aiAvail {
			panic("AI already used or ship does not support AI")
		} else {
			aiAvail = false
		}
	}
}

type VolleyPlanner interface {
	PlanVolleys(ship, target *Ship, theRange int) Volleys
}

type BasicVolleyPlanner struct{}

func (_ BasicVolleyPlanner) PlanVolleys(ship, target *Ship, theRange int) Volleys {
	volleys := make(Volleys, 0, len(ship.Weapons))
	aiAvail := ship.AI
	for weaponIdx, weapon := range ship.Weapons {
		if !weapon.Ranges.IsInRange(theRange) {
			continue
		}
		volley := Volley{
			WeaponIdxs: []int{weaponIdx},
			UseAI:      aiAvail,
		}
		aiAvail = false
		volleys = append(volleys, volley)
	}
	return volleys
}

type PairingVolleyPlanner struct{}

func (_ PairingVolleyPlanner) PlanVolleys(ship, target *Ship, theRange int) Volleys {
	volleys := Volleys{}
	weaponToIdx := map[*Weapon][]int{}
	for weaponIdx, weapon := range ship.Weapons {
		if !weapon.Ranges.IsInRange(theRange) {
			continue
		}
		if !weapon.CanPair {
			volleys = append(volleys, Volley{
				WeaponIdxs: []int{weaponIdx},
			})
			continue
		}
		weaponToIdx[weapon] = append(weaponToIdx[weapon], weaponIdx)
	}
	if len(volleys) == 0 && len(weaponToIdx) == 0 {
		// everything is out of range.
		return nil
	}

	for _, idxSlice := range weaponToIdx {
		var pair []int
		newPair := func() {
			pair = make([]int, 0, 2)
		}
		tryWritePair := func() {
			if len(pair) == 0 {
				return
			}
			volleys = append(volleys, Volley{
				WeaponIdxs: pair,
			})
		}
		for i, weaponIdx := range idxSlice {
			if i%2 == 0 {
				tryWritePair()
				newPair()
			}
			pair = append(pair, weaponIdx)
		}
		tryWritePair()
	}
	if ship.AI {
		// Always put the AI on whatever volley can do the most damage.
		setUseAi := func(idx int, value bool) {
			volleys[idx].UseAI = value
		}
		getDamage := func(volleyIdx int) int {
			return volleys[volleyIdx].Damage(ship)
		}
		maxDamageIdx := 0
		setUseAi(0, true)

		for i := 1; i < len(volleys); i++ {
			volleyDamage := getDamage(i)
			maxDamage := getDamage(maxDamageIdx)
			if volleyDamage > maxDamage {
				setUseAi(maxDamageIdx, false)
				setUseAi(i, true)
				maxDamageIdx = i
			}
		}
	}

	return volleys
}
