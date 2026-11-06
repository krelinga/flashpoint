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
	Hit      bool
}

func (e AttackEvent) Summary() string {
	if e.Hit {
		return fmt.Sprintf("%s hits %s", e.Attacker.Name, e.Target.Name)
	} else {
		return fmt.Sprintf("%s misses %s", e.Attacker.Name, e.Target.Name)
	}
}
