package main

import (
	"cmp"
	"slices"
)

type Salvo struct {
	Target       *Ship
	ThreatValue  int
	ArrivalRound int
}

type SalvoQueue struct {
	salvos []*Salvo
}

func (sq *SalvoQueue) Push(s *Salvo) {
	if s == nil {
		panic("s must not be nil")
	}
	sq.salvos = append(sq.salvos, s)
}

func (sq *SalvoQueue) PopAll(arrivalRound int) []*Salvo {
	slices.SortFunc(sq.salvos, func(a, b *Salvo) int {
		return cmp.Compare(a.ArrivalRound, b.ArrivalRound)
	})
	limit := 0
	for _, s := range sq.salvos {
		if s.ArrivalRound > arrivalRound {
			break
		}
		limit++
	}
	out := sq.salvos[:limit]
	sq.salvos = sq.salvos[limit:]
	return out
}

func SalvoDamage(threat int) int {
	switch {
	case threat <= 1:
		return 3
	case threat <= 3:
		return 5
	case threat <= 6:
		return 7
	case threat <= 10:
		return 10
	default:
		return 11
	}
}
