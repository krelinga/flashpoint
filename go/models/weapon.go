package models

import (
	"errors"
	"fmt"
)

type WeaponData struct {
	Name string
	Ranges Ranges
}

var ErrInvalidWeaponData = errors.New("invalid WeaponData")

type Weapon struct {
	data *WeaponData
}

func NewWeapon(data *WeaponData) (Weapon, error) {
	if err := validateWeaponData(data); err != nil {
		return Weapon{}, err
	}
	return Weapon{data: data}, nil
}

func validateWeaponData(data *WeaponData) error {
	if data == nil {
		return fmt.Errorf("%w: WeaponData is nil", ErrInvalidWeaponData)
	}

	if data.Name == "" {
		return fmt.Errorf("%w: Name must be non-empty", ErrInvalidWeaponData)
	}

	if data.Ranges.IsZero() {
		return fmt.Errorf("%w: Ranges must be non-zero", ErrInvalidWeaponData)
	}

	return nil
}
