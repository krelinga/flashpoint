package tables

import (
	"fmt"
	"reflect"
)

type Table[T any] struct {
	entries []T
	rows, columns int
	lPad, rPad, tPad int
}

func New[T any](rows, columns int) Table[T] {
	if rows <= 0 {
		panic("rows must be positive")
	}
	if columns <= 0 {
		panic("columns must be positive")
	}
	return Table[T]{
		entries: make([]T, rows * columns),
		rows: rows,
		columns: columns,
	}
}

func (t Table[T]) NumCells() int {
	return len(t.entries)
}

func (t Table[T]) Rows() int {
	return t.rows
}

func (t Table[T]) Columns() int {
	return t.columns
}

func (t Table[T]) offset(row, column int) int {
	if row < 0 || row >= t.rows {
		panic("row out of range")
	}
	if column < 0 || column >= t.columns {
		panic("column out of range")
	}
	rowWidth := t.lPad + t.rPad + t.columns
	return rowWidth * (row + t.tPad) + t.lPad + column
}

func (t Table[T]) Get(row, column int) T {
	return t.entries[t.offset(row, column)]
}

func (t Table[T]) Set(row, column int, entry T) {
	t.entries[t.offset(row, column)] = entry
}

type Reducer[T, U any] func(initial U, newValue T) U

func (t Table[T]) ReduceRow[U any](row int, initial U, reducer Reducer[T, U]) U {
	value := initial
	for column := range t.columns {
		value = reducer(value, t.Get(row, column))
	}
	return value
}

func (t Table[T]) ReduceColumn[U any](column int, initial U, reducer Reducer[T, U]) U {
	value := initial
	for row := range t.rows {
		value = reducer(value, t.Get(row, column))
	}
	return value
}

func (t Table[T]) String() string {
	return fmt.Sprintf("table.Table[%s](rows=%d, columns=%d)", reflect.TypeFor[T]().Name(), t.rows, t.columns)
}