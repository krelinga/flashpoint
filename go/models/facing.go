package models

import "iter"

type Facing int

const (
	FacingN Facing = iota + 1
	FacingNE
	FacingSE
	FacingS
	FacingSW
	FacingNW
)

func (f Facing) String() string {
	switch f {
	case FacingN:
		return "N"
	case FacingNE:
		return "NE"
	case FacingSE:
		return "SE"
	case FacingS:
		return "S"
	case FacingSW:
		return "SW"
	case FacingNW:
		return "NW"
	default:
		return "Invalid"
	}
}

func (f Facing) IsValid() bool {
	switch f {
	case FacingN, FacingNE, FacingSE, FacingS, FacingSW, FacingNW:
		return true
	default:
		return false
	}
}

var allFacings = []Facing{FacingN, FacingNE, FacingSE, FacingS, FacingSW, FacingNW}

func NumFacings() int {
	return len(allFacings)
}

func AllFacings() iter.Seq[Facing] {
	return func(yield func(Facing) bool) {
		for _, facing := range allFacings {
			if !yield(facing) {
				return
			}
		}
	}
}
