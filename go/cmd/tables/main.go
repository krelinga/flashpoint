package main

import (
	"fmt"
	"os"

	"github.com/krelinga/flashpoint/go/tables"
)

func main() {
	const rows = 2
	const columns = 3
	value := func(row, column int) string {
		if row < 0 || row > rows {
			panic("row out of range")
		}
		if column < 0 || column > columns {
			panic("column out of range")
		}
		return fmt.Sprintf("row=%d,column=%d", row, column)
	}
	table := tables.New[string](rows, columns)
	for row := range rows {
		for column := range columns {
			table.Set(row, column, value(row, column))
		}
	}

	options := tables.PrettyOptions[string] {
		Title: "table title",
		Box: true,
	}
	if _, err := tables.PrettyWriter(os.Stdout, table, options); err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println(tables.PrettyString(table, options))

	rowHeaders := table.Slice(tables.All(), tables.Count(1))
	fmt.Println(tables.PrettyString(rowHeaders, options))

	colHeaders := table.Slice(tables.Count(1), tables.All())
	fmt.Println(tables.PrettyString(colHeaders, options))

	body := table.Slice(tables.Start(1), tables.Start(1))
	fmt.Println(tables.PrettyString(body, options))
}