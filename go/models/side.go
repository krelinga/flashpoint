package models

import (
	"errors"
	"fmt"
)

type SideData struct {
	Barrier           int
	Armor             int
	MetaMaterialArmor int
}

var ErrInvalidSideData = errors.New("invalid side data")

type Side struct {
	data *SideData
}

func (s Side) IsZero() bool {
	return s.data == nil
}

func NewSide(data *SideData) (Side, error) {
	if err := validateSideData(data); err != nil {
		return Side{}, err
	}
	return Side{data: data}, nil
}

func validateSideData(data *SideData) error {
	if data == nil {
		return fmt.Errorf("%w: data is nil", ErrInvalidSideData)
	}
	if data.Barrier < 0 {
		return fmt.Errorf("%w: barrier is negative", ErrInvalidSideData)
	}
	if data.Armor < 0 {
		return fmt.Errorf("%w: armor is negative", ErrInvalidSideData)
	}
	if data.MetaMaterialArmor < 0 {
		return fmt.Errorf("%w: meta material armor is negative", ErrInvalidSideData)
	}
	return nil
}

func (s Side) Barrier() int {
	return s.data.Barrier
}

func (s Side) Armor() int {
	return s.data.Armor
}

func (s Side) MetaMaterialArmor() int {
	return s.data.MetaMaterialArmor
}
