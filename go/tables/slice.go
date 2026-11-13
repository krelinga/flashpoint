package tables

type Bounds interface {
	getBounds(prefix, size, suffix int) (int, int, int)
}

type startCountBounds struct {
	start, count int
}

func (scb startCountBounds) getBounds(prefix, size, suffix int) (int, int, int) {
	if scb.start < 0 {
		panic("start must be non-negative")
	}
	if scb.count <= 0 {
		panic("count must be positive")
	}
	if scb.start >= size {
		panic("start out of range")
	}
	if scb.start + scb.count >= prefix + size {
		panic("count out of range")
	}
	total := prefix + size + suffix
	newPrefix := prefix + scb.start
	newSize := scb.count
	newSuffix := total - newPrefix - newSize
	return newPrefix, newSize, newSuffix
}

func StartCount(start, count int) Bounds {
	return startCountBounds{
		start: start,
		count: count,
	}
}

type startBounds struct {
	start int
}

func (sb startBounds) getBounds(prefix, size, suffix int) (int, int, int) {
	if sb.start < 0 {
		panic("start must be non-negative")
	}
	if sb.start >= size {
		panic("start out of range")
	}
	newSize := size - sb.start
	newPrefix := prefix + sb.start
	return newPrefix, newSize, suffix
}

func Start(start int) Bounds {
	return startBounds{
		start: start,
	}
}

type countBounds struct {
	count int
}

func (cb countBounds) getBounds(prefix, size, suffix int) (int, int, int) {
	if cb.count < 0 {
		panic("count must be non-negative")
	}
	if cb.count >= size {
		panic("count out of range")
	}
	newSize := cb.count
	newSuffix := suffix + size - cb.count
	return prefix, newSize, newSuffix
}

func Count(count int) Bounds {
	return countBounds{
		count: count,
	}
}

type allBounds struct{}

func (_ allBounds) getBounds(prefix, size, suffix int) (int, int, int) {
	return prefix, size, suffix
}

func All() Bounds {
	return allBounds{}
}

func (t Table[T]) Slice(rows, columns Bounds) Table[T] {
	if rows == nil {
		panic("rows must not be nil")
	}
	if columns == nil {
		panic("columns must not be nil")
	}

	newT := Table[T]{
		entries: t.entries,
	}
	newT.tPad, newT.rows, _ = rows.getBounds(t.tPad, t.rows, 0)
	newT.lPad, newT.columns, newT.rPad = columns.getBounds(t.lPad, t.columns, t.rPad)
	return newT
}