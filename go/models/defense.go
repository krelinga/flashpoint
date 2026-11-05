package models

import (
	"errors"
	"fmt"
)

type DefenseData struct {
	Facing            Facing
	Barrier           int
	Armor             int
	MetaMaterialArmor int
}

var ErrInvalidDefenseData = errors.New("invalid defense data")

type Defense struct {
	data *DefenseData
}

func (s Defense) IsZero() bool {
	return s.data == nil
}

func NewSide(data *DefenseData) (Defense, error) {
	if err := validateSideData(data); err != nil {
		return Defense{}, err
	}
	return Defense{data: data}, nil
}

func validateSideData(data *DefenseData) error {
	if data == nil {
		return fmt.Errorf("%w: data is nil", ErrInvalidDefenseData)
	}
	if !data.Facing.IsValid() {
		return fmt.Errorf("%w: invalid facing", ErrInvalidDefenseData)
	}
	if data.Barrier < 0 {
		return fmt.Errorf("%w: barrier is negative", ErrInvalidDefenseData)
	}
	if data.Armor < 0 {
		return fmt.Errorf("%w: armor is negative", ErrInvalidDefenseData)
	}
	if data.MetaMaterialArmor < 0 {
		return fmt.Errorf("%w: meta material armor is negative", ErrInvalidDefenseData)
	}
	return nil
}

func (s Defense) Barrier() int {
	return s.data.Barrier
}

func (s Defense) Armor() int {
	return s.data.Armor
}

func (s Defense) MetaMaterialArmor() int {
	return s.data.MetaMaterialArmor
}
