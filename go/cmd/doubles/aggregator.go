package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"iter"
	"slices"
	"strings"
)

type KeyValue struct {
	Key, Value string
}

type AggregateKey struct {
	Properties []KeyValue
	Params     []KeyValue
}

func (ak AggregateKey) Equals(other AggregateKey) bool {
	propertiesEqual := slices.Equal(ak.Properties, other.Properties)
	paramsEqual := slices.Equal(ak.Params, other.Params)
	return propertiesEqual && paramsEqual
}

type aggregateCount struct {
	key   AggregateKey
	count int
}

func (ak AggregateKey) String() string {
	var keyString string
	if keyBytes, err := json.Marshal(ak); err != nil {
		panic(err)
	} else {
		keyString = string(keyBytes)
	}
	return keyString
}

type Aggregator struct {
	Properties map[string]func(Result) string
	counts     []*aggregateCount
}

func (a *Aggregator) Total() int {
	var total int
	for _, count := range a.counts {
		total += count.count
	}
	return total
}

func (a *Aggregator) Add(result Result) {
	keyValueCompare := func(a, b KeyValue) int {
		return cmp.Compare(a.Key, b.Key)
	}

	var key AggregateKey
	for name, prop := range a.Properties {
		value := prop(result)
		key.Properties = append(key.Properties, KeyValue{
			Key:   name,
			Value: value,
		})
	}
	slices.SortFunc(key.Properties, keyValueCompare)

	for name, value := range result.ParamValues {
		key.Params = append(key.Params, KeyValue{
			Key:   name,
			Value: value,
		})
	}
	slices.SortFunc(key.Params, keyValueCompare)

	if a.counts == nil {
		a.counts = []*aggregateCount{}
	}
	for _, count := range a.counts {
		if count.key.Equals(key) {
			count.count++
			return
		}
	}
	a.counts = append(a.counts, &aggregateCount{
		key:   key,
		count: 1,
	})
}

func (a *Aggregator) All() iter.Seq2[AggregateKey, int] {
	return func(yield func(AggregateKey, int) bool) {
		for _, count := range a.counts {
			if !yield(count.key, count.count) {
				return
			}
		}
	}
}

func (a *Aggregator) Keys() iter.Seq[AggregateKey] {
	return func(yield func(AggregateKey) bool) {
		for _, count := range a.counts {
			if !yield(count.key) {
				return
			}
		}
	}
}

func (a *Aggregator) String() string {
	if a == nil {
		return "nil"
	}

	compareKeys := func(a, b AggregateKey) int {
		compareKeyValues := func(a, b []KeyValue) int {
			if result := cmp.Compare(len(a), len(b)); result != 0 {
				return result
			}
			for i := range a {
				a := a[i]
				b := b[i]
				if result := cmp.Compare(a.Key, b.Key); result != 0 {
					return result
				}
				if result := cmp.Compare(a.Value, b.Value); result != 0 {
					return result
				}
			}
			return 0
		}
		if result := compareKeyValues(a.Params, b.Params); result != 0 {
			return result
		}
		return compareKeyValues(a.Properties, b.Properties)
	}

	total := a.Total()
	sb := &strings.Builder{}
	sb.WriteString("{\n")
	entries := slices.Clone(a.counts)
	slices.SortFunc(entries, func(a, b *aggregateCount) int {
		return compareKeys(a.key, b.key)
	})
	for _, entry := range entries {
		fraction := float64(entry.count) / float64(total)
		percent := fraction * 100
		fmt.Fprintf(sb, "    %s: %d of %d (%f%%)\n", entry.key.String(), entry.count, total, percent)
	}
	sb.WriteString("}")

	return sb.String()
}
