package models

import (
	"cmp"
	"errors"
	"fmt"
	"iter"
	"slices"
)

var ErrInvalidDefenses = errors.New("invalid defenses")

type Defenses struct {
	data *[]Defense
}

func (d Defenses) IsZero() bool {
	return d.data == nil
}

func NewDefenses(defenses ...Defense) (Defenses, error) {
	data := make([]Defense, len(defenses))
	copy(data, defenses)
	slices.SortFunc(data, func(a, b Defense) int {
		return cmp.Compare(a.Facing(), b.Facing())
	})
	if err := validateDefensesData(data); err != nil {
		return Defenses{}, err
	}
	return Defenses{data: &data}, nil
}

func validateDefensesData(data []Defense) error {
	seen := make(map[Facing]struct{}, NumFacings())
	for _, defense := range data {
		if _, duplicate := seen[defense.Facing()]; duplicate {
			return fmt.Errorf("%w: duplicate defense found for facing %v", ErrInvalidDefenses, defense.Facing())
		}
		seen[defense.Facing()] = struct{}{}
	}
	for facing := range AllFacings() {
		if _, exists := seen[facing]; !exists {
			return fmt.Errorf("%w: missing defense for facing %v", ErrInvalidDefenses, facing)
		}
	}
	return nil
}

func (d Defenses) Values() iter.Seq[Defense] {
	return slices.Values(*d.data)
}

func (d Defenses) At(facing Facing) Defense {
	return (*d.data)[int(facing) - 1]
}