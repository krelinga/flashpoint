package main

type Weapon struct {
	Name   string
	Damage int
	Ranges Ranges
}

var (
	WeaponLaserMk2 = Weapon{
		Name:   "Laser Mk 2",
		Damage: 6,
		Ranges: Ranges{
			Range{Limit: 4, Dice: Multiple(2, Sides(6))},
			Range{Limit: 8, Dice: Sides(10)},
			Range{Limit: 12, Dice: Sides(8)},
			Range{Limit: 16, Dice: Sides(6)},
		},
	}
	WeaponLaserMk3 = Weapon{
		Name:   "Laser Mk 3",
		Damage: 8,
		Ranges: Ranges{
			Range{Limit: 6, Dice: Sides(12)},
			Range{Limit: 12, Dice: Sides(10)},
			Range{Limit: 18, Dice: Sides(8)},
			Range{Limit: 22, Dice: Sides(6)},
		},
	}
	// TODO: support rerolling pulse laser dice.
	WeaponPulseLaserMk2 = Weapon{
		Name: "Pulse Laser Mk 2",
		Damage: 4,
		Ranges: Ranges{
			Range{Limit: 4, Dice: Sides(12)},
			Range{Limit: 8, Dice: Sides(10)},
			Range{Limit: 14, Dice: Sides(8)},
		},
	}
)
