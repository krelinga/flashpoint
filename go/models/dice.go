package models

import (
	"errors"
	"fmt"
)

type DiceData struct {
	Count int
	Sides int
}

var ErrInvalidDiceData = errors.New("invalid DiceData")

type Dice struct {
	data *DiceData
}

func (d Dice) IsZero() bool {
	return d.data == nil
}

func NewDice(data *DiceData) (Dice, error) {
	if err := validateDiceData(data); err != nil {
		return Dice{}, err
	}

	return Dice{data: data}, nil
}

func validateDiceData(data *DiceData) error {
	if data == nil {
		return fmt.Errorf("%w: DiceData is nil", ErrInvalidDiceData)
	}

	if data.Count <= 0 {
		return fmt.Errorf("%w: Count must be positive", ErrInvalidDiceData)
	}

	switch data.Sides {
	case 6, 8, 10, 12:
		// expected valies.
	default:
		return fmt.Errorf("%w: Sides must be in (6, 8, 10, 12)", ErrInvalidDiceData)
	}

	return nil
}

func (d Dice) Count() int {
	return d.data.Count
}

func (d Dice) Sides() int {
	return d.data.Sides
}

var (
	DiceD6  = must(NewDice(&DiceData{Count: 1, Sides: 6}))
	DiceD8  = must(NewDice(&DiceData{Count: 1, Sides: 8}))
	DiceD10 = must(NewDice(&DiceData{Count: 1, Sides: 10}))
	DiceD12 = must(NewDice(&DiceData{Count: 1, Sides: 12}))
	Dice2D6 = must(NewDice(&DiceData{Count: 2, Sides: 6}))
	Dice3D6 = must(NewDice(&DiceData{Count: 3, Sides: 6}))
)
