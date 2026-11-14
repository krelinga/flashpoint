package main

import (
	"fmt"
	"maps"
	"slices"

	"github.com/krelinga/flashpoint/go/tables"
)

type ReportDimension interface {
	findValues(*Aggregator) []string
	match(aggregateKey AggregateKey, value string) bool
}

type ReportParam string

func (r ReportParam) findValues(a *Aggregator) []string {
	found := map[string]struct{}{}
	for key := range a.Keys() {
		for _, paramValue := range key.Params {
			if paramValue.Key == string(r) {
				found[paramValue.Value] = struct{}{}
			}
		}
	}
	out := slices.Collect(maps.Keys(found))
	slices.Sort(out)
	return out
}

func (r ReportParam) match(key AggregateKey, value string) bool {
	for _, param := range key.Params {
		if param.Key == string(r) && param.Value == value {
			return true
		}
	}
	return false
}

type ReportProperty string

func (r ReportProperty) findValues(a *Aggregator) []string {
	found := map[string]struct{}{}
	for key := range a.Keys() {
		for _, propertyValue := range key.Properties {
			if propertyValue.Key == string(r) {
				found[propertyValue.Value] = struct{}{}
			}
		}
	}
	out := slices.Collect(maps.Keys(found))
	slices.Sort(out)
	return out
}

func (r ReportProperty) match(key AggregateKey, value string) bool {
	for _, property := range key.Properties {
		if property.Key == string(r) && property.Value == value {
			return true
		}
	}
	return false
}

type Report struct {
	tables.Table[int]
	rows, columns []string
}


func (r Report) Total() int {
	total := 0
	for row := range r.Rows() {
		for column := range r.Columns() {
			total += r.Get(row, column)
		}
	}
	return total
}

func (r Report) RowTotal(row int) int {
	return r.ReduceRow(row, 0, func(initial, next int) int {
		return initial + next
	})
}

func (r Report) ColumnTotal(column int) int {
	return r.ReduceColumn(column, 0, func(initial, next int) int {
		return initial + next
	})
}

func (r Report) ToTable() tables.Table[string] {
	table := tables.New[string](r.Rows() + 1, r.Columns() + 1)
	for i, row := range r.rows {
		table.Set(i + 1, 0, row)
	}
	for i, column := range r.columns {
		table.Set(0, i + 1, column)
	}
	dataTable := table.Slice(tables.Start(1), tables.Start(1))
	for row := range r.Rows() {
		rowTotal := r.RowTotal(row)
		for column := range r.Columns() {
			columnTotal := r.ColumnTotal(column)
			entry := r.Get(row, column)
			rowPercent := float64(entry) * 100 / float64(rowTotal)
			columnPercent := float64(entry) * 100 / float64(columnTotal)
			dataTable.Set(row, column, fmt.Sprintf("%d (%.02f%% of row, %.02f%% of column)", entry, rowPercent, columnPercent))
		}
	}
	return table
}

// TODO: add NewReportFilter().
func NewReport(a *Aggregator, rows, columns ReportDimension) Report {
	if rows == nil {
		panic("rows must not be nil")
	}
	if columns == nil {
		panic("columns must not be nil")
	}

	rowValues := rows.findValues(a)
	columnValues := columns.findValues(a)
	report := Report{
		Table: tables.New[int](len(rowValues), len(columnValues)),
		rows:    rowValues,
		columns: columnValues,
	}
	for row := range rowValues {
		for column := range columnValues {
			entry := 0
			for key, count := range a.All() {
				if rows.match(key, rowValues[row]) && columns.match(key, columnValues[column]) {
					entry += count
				}
			}
			report.Set(row, column, entry)
		}
	}

	return report
}
