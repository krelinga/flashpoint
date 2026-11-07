package main

import (
	"fmt"
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
	return fmt.Sprintf("%d/%d (%.2f%%)", s.Num, s.Denom, s.Percent())
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