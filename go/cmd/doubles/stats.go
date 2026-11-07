package main

import "fmt"

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