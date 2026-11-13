package main

import (
	"fmt"
	"strings"
)

type Table struct {
	entries       []string
	rows, columns int
}

func (t Table) Rows() int {
	return t.rows
}

func (t Table) Columns() int {
	return t.columns
}

func (t Table) offset(row, column int) int {
	if row < 0 || row > t.Rows() {
		panic("row out of range")
	}
	if column < 0 || column > t.Columns() {
		panic("column out of range")
	}
	return row * (t.Columns() + 1) + column
}

func (t Table) Set(row, column int, value string) {
	t.entries[t.offset(row+1, column+1)] = value
}

func (t Table) SetRowName(row int, value string) {
	t.entries[t.offset(row+1, 0)] = value
}

func (t Table) SetColumnName(column int, value string) {
	t.entries[t.offset(0, column+1)] = value
}

func (t Table) columnWidth(column int) int {
	width := 0
	for row := range t.Rows() + 1 {
		thisWidth := len(t.entries[t.offset(row, column)])
		if thisWidth > width {
			width = thisWidth
		}
	}
	return width
}

func (t Table) String() string {
	if len(t.entries) == 0 {
		return "no data"
	}

	widths := make([]int, t.Columns()+1)
	for column := range t.Columns() + 1 {
		widths[column] = t.columnWidth(column)
	}

	sb := &strings.Builder{}
	for row := range t.Rows() + 1 {
		if row != 0 {
			sb.WriteString("\n")
		}
		for column := range t.Columns() + 1 {
			if column != 0 {
				sb.WriteString(" ")
			}
			formatString := fmt.Sprintf("%%-%ds", widths[column])
			fmt.Fprintf(sb, formatString, t.entries[t.offset(row, column)])
		}
	}
	return sb.String()
}

func NewTable(rows, columns int) Table {
	return Table{
		rows:    rows,
		columns: columns,
		entries: make([]string, (rows+1)*(columns+1)),
	}
}
