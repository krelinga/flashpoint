package main

import (
	"fmt"
	"strings"
)

type Predicate struct {
	Name string
	Test func(Result) bool
}

type WidenedResult struct {
	Result Result
	Matches map[string]struct{}
}

type WidenedResults []WidenedResult

func (wr WidenedResults) Compute(m Matcher) Stat {
	s := Stat{
		Denom: len(wr),
	}
	for _, r := range wr {
		if m.Match(r.Matches) {
			s.Num++
		}
	}
	return s
}

func Widen(results []Result, preds ...Predicate) WidenedResults {
	out := make(WidenedResults, len(results))
	for i, result := range results {
		out[i].Result = result
		for _, pred := range preds {
			if pred.Test(result) {
				if out[i].Matches == nil {
					out[i].Matches = make(map[string]struct{})
				}
				out[i].Matches[pred.Name] = struct{}{}
			}
		}
	}
	return out
}

type Stat struct {
	Num, Denom int
}

func (s Stat) Percent() float64 {
	if s.Denom == 0 {
		return 0
	}
	return float64(s.Num) / float64(s.Denom) * 100
}

func (s Stat) String() string {
	if s.Denom == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%d/%d (%.2f%%)", s.Num, s.Denom, s.Percent())
}

type Matcher interface {
	Match(predNames map[string]struct{}) bool
	String() string
}

func And(matchers ...Matcher) Matcher {
	return andMatcher(matchers)
}

type andMatcher []Matcher

func (a andMatcher) Match(predNames map[string]struct{}) bool {
	for _, m := range a {
		if !m.Match(predNames) {
			return false
		}
	}
	return true
}

func (a andMatcher) String() string {
	sb := &strings.Builder{}
	fmt.Fprint(sb, "And(")
	for i, m := range a {
		if i > 0 {
			fmt.Fprint(sb, ", ")
		}
		fmt.Fprint(sb, m)
	}
	fmt.Fprint(sb, ")")
	return sb.String()
}

func Or(matchers ...Matcher) Matcher {
	return orMatcher(matchers)
}

type orMatcher []Matcher

func (o orMatcher) Match(predNames map[string]struct{}) bool {
	for _, m := range o {
		if m.Match(predNames) {
			return true
		}
	}
	return false
}

func (o orMatcher) String() string {
	sb := &strings.Builder{}
	fmt.Fprint(sb, "Or(")
	for i, m := range o {
		if i > 0 {
			fmt.Fprint(sb, ", ")
		}
		fmt.Fprint(sb, m)
	}
	fmt.Fprint(sb, ")")
	return sb.String()
}

func Not(matcher Matcher) Matcher {
	return notMatcher{matcher}
}

type notMatcher struct {
	child Matcher
}

func (n notMatcher) Match(predNames map[string]struct{}) bool {
	return !n.child.Match(predNames)
}

func (n notMatcher) String() string {
	return fmt.Sprintf("Not(%s)", n.child)
}

func Contains(name string) Matcher {
	return containsMatcher(name)
}

type containsMatcher string

func (c containsMatcher) Match(predNames map[string]struct{}) bool {
	_, found := predNames[string(c)]
	return found
}

func (c containsMatcher) String() string {
	return fmt.Sprintf("Contains(\"%s\")", string(c))
}
