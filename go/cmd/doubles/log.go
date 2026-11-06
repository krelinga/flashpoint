package main

import (
	"fmt"
	"strings"
)

type Logger interface {
	Log(Event)
	Child() Logger
}

type Event interface {
	Summary() string
}

type NullLogger struct{}

func (_ NullLogger) Log(Event) {}

func (_ NullLogger) Child() Logger {
	return NullLogger{}
}

type ConsoleLogger struct {
	depth int
}

func (l ConsoleLogger) Log(e Event) {
	prefix := strings.Repeat("  ", l.depth)
	fmt.Println(prefix + e.Summary())
}

func (l ConsoleLogger) Child() Logger {
	child := ConsoleLogger{depth: l.depth + 1}
	return child
}

type InitiativeEvent struct {
	Winner *Ship
}

func (e InitiativeEvent) Summary() string {
	return fmt.Sprintf("%s wins initiative", e.Winner.Name)
}

type AttackEvent struct {
	Attacker *Ship
	Target   *Ship
	Damage   *Damage
}

func (e AttackEvent) Summary() string {
	if e.Damage == nil {
		return fmt.Sprintf("%s attacks %s and misses", e.Attacker.Name, e.Target.Name)
	}
	return fmt.Sprintf("%s attacks %s and hits: %s", e.Attacker.Name, e.Target.Name, *e.Damage)
}

type NewRoundEvent struct {
	Round int
}

func (e NewRoundEvent) Summary() string {
	return fmt.Sprintf("Round %d begins", e.Round)
}

type NewTurnEvent struct {
	Ship *Ship
}

func (e NewTurnEvent) Summary() string {
	return fmt.Sprintf("%s's turn begins", e.Ship.Name)
}