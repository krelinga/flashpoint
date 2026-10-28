package models

var (
	WeaponLaserMk1 = must(NewWeapon(&WeaponData{
		Name: "Laser Mk 1",
		Ranges: must(NewRanges(RangesData{
			must(NewRange(&RangeData{Min: 0, Max: 4, Dice: Dice3D6})),
			must(NewRange(&RangeData{Min: 5, Max: 8, Dice: Dice2D6})),
			must(NewRange(&RangeData{Min: 9, Max: 10, Dice: DiceD8})),
		})),
	}))
	WeaponLaserMk2 = must(NewWeapon(&WeaponData{
		Name: "Laser Mk 2",
		Ranges: must(NewRanges(RangesData{
			must(NewRange(&RangeData{Min: 0, Max: 4, Dice: Dice2D6})),
			must(NewRange(&RangeData{Min: 5, Max: 8, Dice: DiceD10})),
			must(NewRange(&RangeData{Min: 9, Max: 12, Dice: DiceD8})),
			must(NewRange(&RangeData{Min: 13, Max: 16, Dice: DiceD6})),
		})),
	}))
)
