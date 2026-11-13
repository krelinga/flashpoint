package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"iter"
	"maps"
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

func (ak AggregateKey) String() string {
	var keyString string
	if keyBytes, err := json.Marshal(ak); err != nil {
		panic(err)
	} else {
		keyString = string(keyBytes)
	}
	return keyString
}

func ParseAggregateKeyString(keyString string) AggregateKey {
	buffer := bytes.NewBufferString(keyString)
	keyBytes := buffer.Bytes()
	var key AggregateKey
	if err := json.Unmarshal(keyBytes, &key); err != nil {
		panic(err)
	}
	return key
}

type Aggregator struct {
	Properties map[string]func(Result) string
	counts     map[string]int
}

func (a *Aggregator) Total() int {
	var total int
	for _, count := range a.counts {
		total += count
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
		a.counts = make(map[string]int)
	}
	a.counts[key.String()]++
}

func (a *Aggregator) All() iter.Seq2[AggregateKey, int] {
	return func(yield func(AggregateKey, int) bool) {
		for stringKey, count := range a.counts {
			if !yield(ParseAggregateKeyString(stringKey), count) {
				return
			}
		}
	}
}

func (a *Aggregator) Keys() iter.Seq[AggregateKey] {
	return func(yield func(AggregateKey) bool) {
		for stringKey := range a.counts {
			if !yield(ParseAggregateKeyString(stringKey)) {
				return
			}
		}
	}
}

func (a *Aggregator) String() string {
	if a == nil {
		return "nil"
	}

	total := a.Total()
	sb := &strings.Builder{}
	sb.WriteString("{\n")
	keyStrings := make([]string, 0, len(a.counts))
	keyStrings = slices.AppendSeq(keyStrings, maps.Keys(a.counts))
	slices.Sort(keyStrings)
	for _, keyString := range keyStrings {
		count := a.counts[keyString]
		fraction := float64(count) / float64(total)
		percent := fraction * 100
		fmt.Fprintf(sb, "    %s: %d of %d (%f%%)\n", keyString, count, total, percent)
	}
	sb.WriteString("}")

	return sb.String()
}
