package main

import "github.com/krelinga/flashpoint/go/dice"

type Weapon struct {
	Name   string
	Damage int
	Ranges Ranges
	CanPair bool
}

const PulseTag string = "pulse"

var (
	WeaponLaserMk2 = &Weapon{
		Name:   "Laser Mk 2",
		Damage: 3,
		Ranges: Ranges{
			Range{Limit: 4, Dice: dice.D6().Times(2)},
			Range{Limit: 8, Dice: dice.D10()},
			Range{Limit: 12, Dice: dice.D8()},
			Range{Limit: 16, Dice: dice.D6()},
		},
		CanPair: true,
	}
	WeaponLaserMk3 = &Weapon{
		Name:   "Laser Mk 3",
		Damage: 4,
		Ranges: Ranges{
			Range{Limit: 6, Dice: dice.D12()},
			Range{Limit: 12, Dice: dice.D10()},
			Range{Limit: 18, Dice: dice.D8()},
			Range{Limit: 22, Dice: dice.D6()},
		},
		CanPair: true,
	}
	WeaponPulseLaserMk2 = &Weapon{
		Name:   "Pulse Laser Mk 2",
		Damage: 4,
		Ranges: Ranges{
			Range{Limit: 4, Dice: dice.D12(PulseTag)},
			Range{Limit: 8, Dice: dice.D10(PulseTag)},
			Range{Limit: 14, Dice: dice.D8(PulseTag)},
		},
	}
)
