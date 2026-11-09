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
	Attacker   *Ship
	Target     *Ship
	Range      int
	Weapon     Weapon
	Roll       int
	AnyRerolls bool
	Damage     Damage
	AI         bool
}

func (e AttackEvent) Summary() string {
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "%s attacks %s", e.Attacker.Name, e.Target.Name)
	modifiers := e.AI || e.AnyRerolls
	if modifiers {
		sb.WriteString(" (")
		var needSep bool
		writeModifier := func(modifier string) {
			if needSep {
				sb.WriteString(", ")
			}
			sb.WriteString(modifier)
			needSep = true
		}
		if e.AI {
			writeModifier("with AI")
		}
		if e.AnyRerolls {
			writeModifier("with rerolls")
		}
		sb.WriteString(")")
	}
	fmt.Fprintf(sb, " at range %d with %s and", e.Range, e.Weapon.Name)
	hit := e.Roll >= e.Target.Aspect
	if hit {
		fmt.Fprint(sb, " hits")
	} else {
		fmt.Fprint(sb, " misses")
	}
	fmt.Fprintf(sb, " (roll: %d vs aspect: %d)", e.Roll, e.Target.Aspect)
	if hit {
		fmt.Fprintf(sb, " and deals damage: %s", e.Damage)
	}
	return sb.String()
}

type NewRoundEvent struct {
	Round int
}

func (e NewRoundEvent) Summary() string {
	return fmt.Sprintf("Round %d begins", e.Round)
}

type WeaponAttackTurnEvent struct {
	Ship *Ship
}

func (e WeaponAttackTurnEvent) Summary() string {
	return fmt.Sprintf("%s's weapon attack turn begins", e.Ship.Name)
}

type PDCTurnEvent struct {
	Target     *Ship
	BaseThreat int
	PdcHits    int
	Damage     Damage
}

func (e PDCTurnEvent) AdjustedThreat() int {
	return max(e.BaseThreat-e.PdcHits, 0)
}

func (e PDCTurnEvent) Summary() string {
	sb := &strings.Builder{}
	fmt.Fprintf(sb, "%s conducts PDC: %d threat ", e.Target.Name, e.BaseThreat)
	adjustedThreat := e.AdjustedThreat()
	if adjustedThreat == e.BaseThreat {
		fmt.Fprint(sb, "(PDC ineffective)")
	} else {
		fmt.Fprintf(sb, "(reduced to %d by %d PDC hits)", adjustedThreat, e.PdcHits)
	}
	if e.AdjustedThreat() > 0 {
		fmt.Fprintf(sb, ", damage was %s", e.Damage)
	}
	return sb.String()
}

type SalvoLaunchEvent struct {
	Attacker, Target *Ship
	NumSalvos        int
	SalvoThreat      int
}

func (e SalvoLaunchEvent) Summary() string {
	return fmt.Sprintf("%s launches %d salvos of threat %d at %s", e.Attacker.Name, e.NumSalvos, e.SalvoThreat, e.Target.Name)
}
