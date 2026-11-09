package dice

import (
	"iter"
	"slices"
)

type Model struct {
	sides int
	tags  []string
}

func (m Model) IsZero() bool {
	return m.sides == 0 && m.tags == nil
}

func (m Model) panicIfZero() {
	if m.IsZero() {
		panic("Non-Zero value of dice.Model is required")
	}
}

func (m Model) Sides() int {
	m.panicIfZero()
	return m.sides
}

func (m Model) HasTag(tag string) bool {
	m.panicIfZero()
	return slices.Contains(m.tags, tag)
}

func (m Model) Find(tag string) iter.Seq[Model] {
	return func(yield func(Model) bool) {
		if m.HasTag(tag) {
			yield(m)
		}
	}
}

func (m Model) roll(roller Roller) int {
	m.panicIfZero()
	return roller.Roll(m.sides)
}

func (m Model) rollAndAppend(roller Roller, results *ResultPool) {
	m.panicIfZero()
	*results = append(*results, Result{
		Model: m,
		Value: m.roll(roller),
	})
}

func (m Model) count() int {
	m.panicIfZero()
	return 1
}

func (m Model) Times(times int) ModelPool {
	if times <= 0 {
		panic("times must be positive")
	}
	return timesPool{
		times: times,
		model: m,
	}
}

type timesPool struct {
	times int
	model Model
}

func (tp timesPool) Find(tag string) iter.Seq[Model] {
	return func(yield func(Model) bool) {
		if !tp.model.HasTag(tag) {
			return
		}
		for range tp.times {
			if !yield(tp.model) {
				return
			}
		}
	}
}

func (tp timesPool) count() int {
	return tp.times
}

func (tp timesPool) rollAndAppend(roller Roller, results *ResultPool) {
	for range tp.times {
		tp.model.rollAndAppend(roller, results)
	}
}

func newModel(sides int, tags []string) Model {
	return Model{
		sides: sides,
		tags:  slices.Clone(tags),
	}
}

func D6(tags ...string) Model {
	return newModel(6, tags)
}

func D8(tags ...string) Model {
	return newModel(8, tags)
}

func D10(tags ...string) Model {
	return newModel(10, tags)
}

func D12(tags ...string) Model {
	return newModel(10, tags)
}

type ModelPool interface {
	Find(tag string) iter.Seq[Model]
	rollAndAppend(Roller, *ResultPool)
	count() int
}

func Add(pools ...ModelPool) ModelPool {
	if len(pools) == 0 {
		panic("pools must not be empty")
	}
	return addedModelPool(pools)
}

type addedModelPool []ModelPool

func (a addedModelPool) Find(tag string) iter.Seq[Model] {
	return func(yield func(Model) bool) {
		for _, entry := range a {
			for pool := range entry.Find(tag) {
				if !yield(pool) {
					return
				}
			}
		}
	}
}

func (a addedModelPool) rollAndAppend(roller Roller, results *ResultPool) {
	for _, entry := range a {
		entry.rollAndAppend(roller, results)
	}
}

func (a addedModelPool) count() int {
	var size int
	for _, entry := range a {
		size += entry.count()
	}
	return size
}
