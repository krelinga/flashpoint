package models

type Faction int

const (
	FactionTerranSovereignty Faction = iota + 1
)

func (f Faction) String() string {
	switch f {
	case FactionTerranSovereignty:
		return "Terran Sovereignty"
	default:
		return "Invalid"
	}
}

func (f Faction) IsValid() bool {
	switch f {
	case FactionTerranSovereignty:
		return true
	default:
		return false
	}
}
