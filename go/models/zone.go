package models

import (
	"iter"
	"slices"
)

type Zone int

const (
	ZoneFore Zone = iota + 1
	ZonePort
	ZoneAft
	ZoneStarboard
	ZoneCore
)

func (z Zone) String() string {
	switch z {
	case ZoneFore:
		return "Fore"
	case ZonePort:
		return "Port"
	case ZoneAft:
		return "Aft"
	case ZoneStarboard:
		return "Starboard"
	case ZoneCore:
		return "Core"
	default:
		return "Unknown"
	}
}

func (z Zone) IsValid() bool {
	switch z {
	case ZoneFore, ZonePort, ZoneAft, ZoneStarboard, ZoneCore:
		return true
	default:
		return false
	}
}

var allZones = []Zone{ZoneFore, ZonePort, ZoneAft, ZoneStarboard, ZoneCore}

func NumZones() int {
	return len(allZones)
}

func AllZones() iter.Seq[Zone] {
	return slices.Values(allZones)
}