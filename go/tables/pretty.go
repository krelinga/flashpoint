package tables

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type PrettyOptions[T any] struct {
	Title                   string
	Box                     bool
	HeaderRow, HeaderColumn bool
	StringFunc              func(T) string
	UTF8                    bool
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

	type charSetType struct {
		boxSide, boxTop                                    string
		boxCornerTL, boxCornerTR, boxCornerBL, boxCornerBR string
		titleUnderline                                     string
		headerUnderline, headerSep, headerIntersection     string
		cellSep                                            string
		boxPadding                                         string
	}

	var charSet charSetType
	if options.UTF8 {
		charSet.boxSide = "│"
		charSet.boxTop = "─"
		charSet.boxCornerTL = "┌"
		charSet.boxCornerTR = "┐"
		charSet.boxCornerBL = "└"
		charSet.boxCornerBR = "┘"
		charSet.titleUnderline = "═"
		charSet.headerUnderline = "─"
		charSet.headerSep = "│"
		charSet.headerIntersection = "┼"
		charSet.cellSep = " "
		charSet.boxPadding = " "
	} else {
		charSet.boxSide = "|"
		charSet.boxTop = "-"
		charSet.boxCornerTL = "+"
		charSet.boxCornerTR = "+"
		charSet.boxCornerBL = "+"
		charSet.boxCornerBR = "+"
		charSet.titleUnderline = "="
		charSet.headerUnderline = "-"
		charSet.headerSep = "|"
		charSet.headerIntersection = "+"
		charSet.cellSep = " "
		charSet.boxPadding = " "
	}
	const cellPadding = 1
	const boxPadding = 1

	dataWidth := func() int {
		total := 0
		for column := range table.Columns() {
			var sep string
			if column == 1 {
				sep = charSet.headerSep
			} else {
				sep = charSet.cellSep
			}
			if column != 0 {
				total += cellPadding*2 + utf8.RuneCountInString(sep)
			}
			total += columnChars[column]
		}
		return max(total, len(options.Title))
	}()

	print := func(args ...any) (newChars int) {
		for _, arg := range args {
			if err != nil {
				return
			}
			var newN int
			newString := fmt.Sprint(arg)
			newChars += utf8.RuneCountInString(newString)
			newN, err = io.WriteString(writer, newString)
			n += newN
		}
		return
	}

	printLn := func(args ...any) {
		if options.Box {
			print(charSet.boxSide)
			print(strings.Repeat(charSet.boxPadding, boxPadding))
		}

		newChars := print(args...)

		if options.Box {
			print(strings.Repeat(" ", dataWidth-newChars))
			print(strings.Repeat(charSet.boxPadding, boxPadding))
			print(charSet.boxSide)
		}

		print("\n")
	}

	if options.Box {
		// Top line
		print(
			charSet.boxCornerTL,
			strings.Repeat(charSet.boxTop, dataWidth+2*boxPadding),
			charSet.boxCornerTR,
			"\n",
		)
	}

	if options.Title != "" {
		printLn(options.Title)
		printLn(strings.Repeat(charSet.titleUnderline, utf8.RuneCountInString(options.Title)))
	}

	for row := range table.Rows() {
		if row == 1 && options.HeaderRow {
			left := strings.Repeat(charSet.headerUnderline, columnChars[0] + cellPadding)
			pivot := charSet.headerUnderline
			if options.HeaderColumn {
				pivot = charSet.headerIntersection
			}
			right := strings.Repeat(charSet.headerUnderline, dataWidth - utf8.RuneCountInString(left) - utf8.RuneCountInString(pivot))
			printLn(left, pivot, right)
		}

		cells := make([]any, 0, table.Columns()*5)
		for column := range table.Columns() {
			if column != 0 {
				padding := strings.Repeat(" ", cellPadding)
				sep := charSet.cellSep
				if column == 1 && options.HeaderColumn {
					sep = charSet.headerSep
				}
				cells = append(cells, padding, sep, padding)
			}
			cellString := options.StringFunc(table.Get(row, column))
			neededPadding := columnChars[column] - utf8.RuneCountInString(cellString)
			cells = append(cells, cellString, strings.Repeat(" ", neededPadding))
		}
		printLn(cells...)
	}

	if options.Box {
		// Bottom line
		print(
			charSet.boxCornerBL,
			strings.Repeat(charSet.boxTop, dataWidth+2*boxPadding),
			charSet.boxCornerBR,
		)
	}

	/*const horizontalDivider = " | "
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
	boxTop := func() {
		if options.Box {
			if !options.UTF8 {
				writeF(" %s\n", strings.Repeat("-", numChars))
				return
			}

			writeF("┌%s┐\n", strings.Repeat("─", numChars))
		}
	}
	boxBottom := func() {
		if options.Box {
			if !options.UTF8 {
				writeF("\n %s", strings.Repeat("-", numChars))
				return
			}
			writeF("└%s┘\n", strings.Repeat("─", numChars))
		}
	}
	boxVerticalLine := func() {
		if options.Box {
			if !options.UTF8 {
				write("|")
			} else {
				write("│")
			}
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

	boxTop()
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
		if row == 1 && options.HeaderRow {
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
			if column == 1 && options.HeaderColumn {
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
	boxBottom()*/

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
