package models

import (
	"errors"
	"fmt"
	"iter"
	"slices"
)

type RangesData []Range

var ErrInvalidRangesData = errors.New("invalid RangesData")

type Ranges struct {
	data *RangesData
}

func (r Ranges) IsZero() bool {
	return r.data == nil
}

func NewRanges(data RangesData) (Ranges, error) {
	if err := validateRangesData(data); err != nil {
		return Ranges{}, err
	}

	return Ranges{data: &data}, nil
}

func validateRangesData(data RangesData) error {
	if data == nil {
		return fmt.Errorf("%w: RangesData is nil", ErrInvalidRangesData)
	}

	if len(data) == 0 {
		return fmt.Errorf("%w: zero length RangesData", ErrInvalidRangesData)
	}

	for i, r := range data {
		if r.IsZero() {
			return fmt.Errorf("%w: RangeData[%d] is zero", ErrInvalidRangesData, i)
		}
	}

	for i := 1; i < len(data); i++ {
		if data[i].Min() < data[i-1].Max() {
			return fmt.Errorf("%w: RangeData[%d] out of order", ErrInvalidRangesData, i)
		}
	}

	return nil
}

func (r Ranges) All() iter.Seq2[int, Range] {
	return slices.All(*r.data)
}

func (r Ranges) Values() iter.Seq[Range] {
	return slices.Values(*r.data)
}

func (r Ranges) Len() int {
	return len(*r.data)
}

func (r Ranges) At(index int) Range {
	return (*r.data)[index]
}
