package models

import (
	"errors"
	"fmt"
)

type ArcData struct {
	Start int
	Width int
	Fanout bool
}

var ErrInvalidArcData = errors.New("invalid ArcData")

type Arc struct {
	data *ArcData
}

func (a Arc) IsZero() bool {
	return a.data == nil
}

func NewArc(data *ArcData) (Arc, error) {
	if err := validateArcData(data); err != nil {
		return Arc{}, err
	}

	return Arc{data: data}, nil
}

func validateArcData(data *ArcData) error {
	if data == nil {
		return fmt.Errorf("%w: ArcData is nil", ErrInvalidArcData)
	}

	if data.Start >= -179 && data.Start <= 180 {
		// Allowed values, no error
	} else {
		return fmt.Errorf("%w: Start out of range, expeced -179 <= Start <= 180", ErrInvalidArcData)
	}

	if data.Width >= 0 && data.Width <= 360 {
		// Allowed values, no error
	} else {
		return fmt.Errorf("%w: Width out of range, expected 0 <= Width <= 360", ErrInvalidArcData)
	}

	if data.Fanout {
		switch data.Start {
		case -180, -120, -60, 0, 60, 120, 180:
			// Allowed values, no error.
		default:
			return fmt.Errorf("%w: when Fanout == true the only allowed Start values are -180, -120, -60, 0, 60, 120, 180", ErrInvalidArcData)
		}

		if data.Width != 0 {
			return fmt.Errorf("%w: when Fanout == true the only allowed Width value is 0", ErrInvalidArcData)
		}
	}

	return nil
}

func (a Arc) Start() int {
	return a.data.Start
}

func (a Arc) Width() int {
	return a.data.Width
}

func (a Arc) Fanout() bool {
	return a.data.Fanout
}