package models

import (
	"errors"
	"fmt"
)

type AmmoData struct {
	Name string
	Speed int
	Range int
}

var ErrInvalidAmmoData = errors.New("invalid AmmoData")

type Ammo struct {
	data *AmmoData
}

func (a Ammo) IsZero() bool {
	return a.data == nil
}

func NewAmmo(data *AmmoData) (Ammo, error) {
	if err := validateAmmoData(data); err != nil {
		return Ammo{}, err
	}

	return Ammo{data: data}, nil
}

func validateAmmoData(data *AmmoData) error {
	if data == nil {
		return fmt.Errorf("%w: AmmoData is nil", ErrInvalidAmmoData)
	}

	if data.Name == "" {
		return fmt.Errorf("%w: Name must be non-empty", ErrInvalidAmmoData)
	}

	if data.Speed <= 0 {
		return fmt.Errorf("%w: Speed must be positive")
	}

	if data.Range <= 0 {
		return fmt.Errorf("%w: Range must be positive")
	}

	return nil
}

func (a Ammo) Name() string {
	return a.data.Name
}

func (a Ammo) Speed() int {
	return a.data.Speed
}

func (a Ammo) Range() int {
	return a.data.Range
}