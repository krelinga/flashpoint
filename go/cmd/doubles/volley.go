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

type BasicVolleyPlanner struct {}

func (_ BasicVolleyPlanner) PlanVolleys(ship, target *Ship, theRange int) Volleys {
	volleys := make(Volleys, 0, len(ship.Weapons))
	aiAvail := ship.AI
	for weaponIdx, weapon := range ship.Weapons {
		if !weapon.Ranges.IsInRange(theRange) {
			continue
		}
		volley := Volley{
			WeaponIdxs: []int{weaponIdx},
			UseAI: aiAvail,
		}
		aiAvail = false
		volleys = append(volleys, volley)
	}
	return volleys
}