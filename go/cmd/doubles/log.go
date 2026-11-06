package main

import "fmt"

type Logger interface {
	Log(Event)
}

type Event interface {
	Summary() string
}

type NullLogger struct{}

func (_ NullLogger) Log(Event) {}

type ConsoleLogger struct{}

func (_ ConsoleLogger) Log(e Event) {
	fmt.Println(e.Summary())
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
