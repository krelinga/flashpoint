package models

import (
	"errors"
	"fmt"
)

type RangeData struct {
	Min, Max int
	Dice     Dice
}

var ErrInvalidRangeData = errors.New("invalid RangeData")

type Range struct {
	data *RangeData
}

func (r Range) IsZero() bool {
	return r.data == nil
}

func NewRange(data *RangeData) (Range, error) {
	if err := validateRangeData(data); err != nil {
		return Range{}, nil
	}

	return Range{data: data}, nil
}

func validateRangeData(data *RangeData) error {
	if data == nil {
		return fmt.Errorf("%w: RangeData is nil", ErrInvalidRangeData)
	}

	if data.Min <= 0 {
		return fmt.Errorf("%w: Min must be positive", ErrInvalidRangeData)
	}

	if data.Max <= data.Min {
		return fmt.Errorf("%w: Max must be greater than Min", ErrInvalidRangeData)
	}

	if data.Dice.IsZero() {
		return fmt.Errorf("%w: Dice is zero", ErrInvalidRangeData)
	}

	return nil
}

func (r Range) Min() int {
	return r.data.Min
}

func (r Range) Max() int {
	return r.data.Max
}

func (r Range) Dice() Dice {
	return r.data.Dice
}
