package models

type ShipData struct {
	Faction Faction
	Name string
	SizeClass float64
	Defenses Defenses
	Mounts []Mount
}

type Ship struct {}