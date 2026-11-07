package main

import (
	"fmt"
	"slices"
	"strings"
)

type Stat struct {
	Matcher Matcher
	Num, Denom int
}

func (s Stat) Percent() float64 {
	if s.Denom == 0 {
		return 0
	}
	return float64(s.Num) / float64(s.Denom) * 100
}

func (s Stat) FractionString() string {
	if s.Denom == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%s/%s (%.2f%%)", prettyInt(s.Num), prettyInt(s.Denom), s.Percent())
}

func (s Stat) String() string {
	return fmt.Sprintf("%s %s", s.Matcher, s.FractionString())
}

type StatTable []Stat

func (st StatTable) String() string {
	const (
		matcherTitle = "Matcher"
		statTitle = "Stat"
	)
	longestName := len(matcherTitle)
	longestFraction := len(statTitle)
	for _, s := range st {
		if nameLen := len(s.Matcher.String()); nameLen > longestName {
			longestName = nameLen
		}
		if fractionLen := len(s.FractionString()); fractionLen > longestFraction {
			longestFraction = fractionLen
		}
	}
	fmtString := fmt.Sprintf("%%-%ds | %%%ds", longestName, longestFraction)
	sb := &strings.Builder{}
	fmt.Fprintf(sb, fmtString, matcherTitle, statTitle)
	fmt.Fprint(sb, "\n")
	fmt.Fprintf(sb, fmtString, strings.Repeat("=", longestName), strings.Repeat("=", longestFraction))
	for _, s := range st {
		fmt.Fprint(sb, "\n")
		fmt.Fprintf(sb, fmtString, s.Matcher.String(), s.FractionString())
	}
	return sb.String()
}

func prettyInt(raw int) string {
	if raw == 0 {
		return "0"
	} else if raw < 0 {
		panic("negative numbers are not supported")
	}
	var parts []string
	var cols int
	for ; raw > 0; raw /= 10 {
		if cols > 0 && cols % 3 == 0 {
			parts = append(parts, ",")
		}
		parts = append(parts, fmt.Sprintf("%d", raw % 10))
		cols++
	}
	slices.Reverse(parts)
	sb := &strings.Builder{}
	for _, part := range parts {
		sb.WriteString(part)
	}
	return sb.String()
}