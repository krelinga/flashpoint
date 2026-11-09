package dice

import "iter"

type Result struct {
	Model      Model
	Value      int
	NumRerolls int
}

func (r Result) Reroll(roller Roller) Result {
	return Result{
		Model:      r.Model,
		Value:      max(r.Value, r.Model.roll(roller)),
		NumRerolls: r.NumRerolls + 1,
	}
}

type ResultPool []Result

func (rp ResultPool) Find(tag string) iter.Seq2[int, Result] {
	return func(yield func(int, Result) bool) {
		for i, r := range rp {
			if !r.Model.HasTag(tag) {
				continue
			}
			if !yield(i, r) {
				return
			}
		}
	}
}

func (rp ResultPool) Total() int {
	var total int
	for _, r := range rp {
		total += r.Value
	}
	return total
}

func Roll(models ModelPool, roller Roller) ResultPool {
	if models == nil {
		panic("models must not be nil")
	}
	if roller == nil {
		panic("roller must not be nil")
	}
	results := make(ResultPool, 0, models.count())
	models.rollAndAppend(roller, &results)
	return results
}
