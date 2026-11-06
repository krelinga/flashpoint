package main

import "math/rand/v2"

type Roller interface {
	Roll(sides int) int
}

type RandomRoller struct {
	Rand *rand.Rand
}

func (r *RandomRoller) Roll(sides int) int {
	return r.Rand.IntN(sides) + 1
}

type Dice interface {
	Roll(roller Roller) int
}

type multipleDice struct {
	Count int
	Die Dice
}

func (m *multipleDice) Roll(roller Roller) int {
	total := 0
	for i := 0; i < m.Count; i++ {
		total += m.Die.Roll(roller)
	}
	return total
}

func Multiple(count int, die Dice) Dice {
	return &multipleDice{Count: count, Die: die}
}

type addDice []Dice

func (a addDice) Roll(roller Roller) int {
	total := 0
	for _, die := range a {
		total += die.Roll(roller)
	}
	return total
}

func Add(dice ...Dice) Dice {
	return addDice(dice)
}

type sidesDice int

func (s sidesDice) Roll(roller Roller) int {
	return roller.Roll(int(s))
}

func Sides(sides int) Dice {
	return sidesDice(sides)
}