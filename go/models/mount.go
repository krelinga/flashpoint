package models

import (
	"errors"
	"fmt"
)

type MountData struct {
	Count  int
	Weapon Weapon
	Arc    Arc
}

var ErrInvalidMountData = errors.New("invalid mount data")

type Mount struct {
	data *MountData
}

func (m Mount) IsZero() bool {
	return m.data == nil
}

func NewMount(data *MountData) (Mount, error) {
	if err := validateMountData(data); err != nil {
		return Mount{}, err
	}
	return Mount{data: data}, nil
}

func validateMountData(data *MountData) error {
	if data == nil {
		return fmt.Errorf("%w: data is nil", ErrInvalidMountData)
	}
	if data.Count <= 0 {
		return fmt.Errorf("%w: count is not positive", ErrInvalidMountData)
	}
	if data.Weapon.IsZero() {
		return fmt.Errorf("%w: weapon is zero", ErrInvalidMountData)
	}
	if data.Arc.IsZero() {
		return fmt.Errorf("%w: arc is zero", ErrInvalidMountData)
	}
	return nil
}

func (m Mount) Count() int {
	return m.data.Count
}

func (m Mount) Weapon() Weapon {
	return m.data.Weapon
}

func (m Mount) Arc() Arc {
	return m.data.Arc
}
