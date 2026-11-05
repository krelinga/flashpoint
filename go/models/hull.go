package models

import (
	"errors"
	"fmt"
	"iter"
)

type HullBox int

func (h HullBox) String() string {
	switch h {
	case HullBoxRegular:
		return "Regular"
	case HullBoxStress:
		return "Stress"
	case HullBoxCritical:
		return "Critical"
	default:
		return "Invalid"
	}
}

func (h HullBox) IsValid() bool {
	switch h {
	case HullBoxRegular, HullBoxStress, HullBoxCritical:
		return true
	default:
		return false
	}
}

const (
	HullBoxRegular HullBox = iota + 1
	HullBoxStress
	HullBoxCritical
)

type hullBoxesData struct {
	box   HullBox
	count int
}

func NewHullBoxes(box HullBox, count int) (HullBoxes, error) {
	data := &hullBoxesData{
		box:   box,
		count: count,
	}
	if err := validateHullBoxSpanData(data); err != nil {
		return HullBoxes{}, err
	}
	return HullBoxes{
		data: data,
	}, nil
}

func validateHullBoxSpanData(data *hullBoxesData) error {
	if !data.box.IsValid() {
		return fmt.Errorf("invalid hull box: %v", data.box)
	}
	if data.count < 0 {
		return fmt.Errorf("invalid count: %d", data.count)
	}
	return nil
}

var ErrInvalidHullBoxes = errors.New("invalid hull boxes")

type HullBoxes struct {
	data *hullBoxesData
}

func (h HullBoxes) IsZero() bool {
	return h.data == nil
}

func (h HullBoxes) Box() HullBox {
	return h.data.box
}

func (h HullBoxes) Count() int {
	return h.data.count
}

type hullData struct {
	zone Zone
	boxes []HullBoxes
}

type Hull struct {
	data *hullData
}

func (h Hull) IsZero() bool {
	return h.data == nil
}

func NewHull(zone Zone, boxes ...HullBoxes) (Hull, error) {
	data := &hullData{
		zone: zone,
		boxes: make([]HullBoxes, len(boxes)),
	}
	copy(data.boxes, boxes)
	if err := validateHullData(data); err != nil {
		return Hull{}, err
	}
	return Hull{data: data}, nil
}

func validateHullData(data *hullData) error {
	if !data.zone.IsValid() {
		return fmt.Errorf("%w: invalid zone: %v", ErrInvalidHullBoxes, data.zone)
	}
	if len(data.boxes) == 0 {
		return fmt.Errorf("%w: no hull boxes provided", ErrInvalidHullBoxes)
	}
	for i, b := range data.boxes {
		if b.IsZero() {
			return fmt.Errorf("%w: invalid hull box at index %d", ErrInvalidHullBoxes, i)
		}
	}
	return nil
}

func (h Hull) Zone() Zone {
	return h.data.zone
}

func (h Hull) All() iter.Seq2[int, HullBox] {
	return func(yield func(int, HullBox) bool) {
		i := 0
		for _, boxes := range h.data.boxes {
			for range boxes.Count() {
				if !yield(i, boxes.Box()) {
					return
				}
				i++
			}
		}
	}
}

func (h Hull) Total() int {
	total := 0
	for _, boxes := range h.data.boxes {
		total += boxes.Count()
	}
	return total
}

func (h Hull) Filter(hullBox HullBox) iter.Seq[int] {
	return func(yield func(int) bool) {
		i := 0
		for _, boxes := range h.data.boxes {
			if boxes.Box() != hullBox {
				i += boxes.Count()
				continue
			}

			for range boxes.Count() {
				if !yield(i) {
					return
				}
				i++
			}
		}
	}
}
