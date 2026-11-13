package tables

import (
	"fmt"
	"io"
	"strings"
)

type PrettyOptions[T any] struct {
	Title      string
	Box        bool
	StringFunc func(T) string
}

func PrettyWriter[T any](writer io.Writer, table Table[T], options PrettyOptions[T]) (n int, err error) {
	if options.StringFunc == nil {
		options.StringFunc = func(raw T) string {
			return fmt.Sprint(raw)
		}
	}
	if table.NumCells() == 0 {
		return fmt.Fprint(writer, "no data")
	}
	columnChars := make([]int, table.Columns())
	for column := range table.Columns() {
		longest := table.ReduceColumn(column, 0, func(initial int, next T) int {
			nextLen := len(options.StringFunc(next))
			return max(initial, nextLen)
		})
		columnChars[column] = longest
	}
	const horizontalDivider = " | "
	const horizontalSpace = "   "
	numChars := func() int {
		tableCols := 0
		for column := range table.Columns() {
			if column == 1 {
				tableCols += len(horizontalDivider)
			} else if column != 0 {
				tableCols += len(horizontalSpace)
			}
			tableCols += columnChars[column]
		}
		return max(len(options.Title), tableCols)
	}()

	write := func(parts ...any) {
		if err != nil {
			return
		}
		var newN int
		newN, err = fmt.Fprint(writer, parts...)
		n += newN
	}
	writeF := func(format string, parts ...any) {
		if err != nil {
			return
		}
		var newN int
		newN, err = fmt.Fprintf(writer, format, parts...)
		n += newN
	}
	newLine := func() {
		write("\n")
	}
	boxHorizontalLine := func() {
		if options.Box {
			writeF(" %s", strings.Repeat("-", numChars))
		}
	}
	boxVerticalLine := func() {
		if options.Box {
			write("|")
		}
	}
	colDivider := func() {
		write(horizontalDivider)
	}
	colSpace := func() {
		write(horizontalSpace)
	}
	rowDivider := func(column int) {
		write(strings.Repeat("-", columnChars[column]))
	}

	boxHorizontalLine()
	if options.Box {
		newLine()
	}
	if options.Title != "" {
		boxVerticalLine()
		format := fmt.Sprintf("%%-%ds", numChars)
		writeF(format, options.Title)
		boxVerticalLine()
		newLine()

		boxVerticalLine()
		writeF(format, strings.Repeat("=", len(options.Title)))
		boxVerticalLine()
		newLine()

		boxVerticalLine()
		writeF(format, "")
		boxVerticalLine()
		newLine()
	}

	for row := range table.Rows() {
		if row > 0 {
			newLine()
		}
		boxVerticalLine()
		if row == 1 {
			for column := range table.Columns() {
				if column > 0 {
					colSpace()
				}
				rowDivider(column)
			}
			boxVerticalLine()
			newLine()
			boxVerticalLine()
		}
		for column := range table.Columns() {
			if column == 1 {
				colDivider()
			} else if column != 0 {
				colSpace()
			}
			format := fmt.Sprintf("%%-%ds", columnChars[column])
			writeF(format, options.StringFunc(table.Get(row, column)))
		}
		boxVerticalLine()
	}
	if options.Box {
		newLine()
	}
	boxHorizontalLine()

	return
}

func PrettyString[T any](table Table[T], options PrettyOptions[T]) string {
	sb := &strings.Builder{}
	_, err := PrettyWriter(sb, table, options)
	if err != nil {
		panic(err)
	}
	return sb.String()
}
