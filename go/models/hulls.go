package models

import (
	"cmp"
	"errors"
	"fmt"
	"iter"
	"slices"
)

type Hulls struct {
	data *[]Hull
}

var ErrInvalidHulls = errors.New("invalid hulls")

func (h Hulls) IsZero() bool {
	return h.data == nil
}

func NewHulls(hulls ...Hull) (Hulls, error) {
	data := make([]Hull, len(hulls))
	copy(data, hulls)
	slices.SortFunc(data, func(a, b Hull) int {
		return cmp.Compare(a.Zone(), b.Zone())
	})
	if err := validateHullsData(data); err != nil {
		return Hulls{}, err
	}
	return Hulls{data: &data}, nil
}

func validateHullsData(data []Hull) error {
	seen := make(map[Zone]struct{}, NumZones())
	for _, hull := range data {
		zone := hull.Zone()
		if _, exists := seen[zone]; exists {
			return fmt.Errorf("%w: duplicate zone: %v", ErrInvalidHulls, zone)
		}
		seen[zone] = struct{}{}
	}
	for zone := range AllZones() {
		if _, exists := seen[zone]; !exists {
			return fmt.Errorf("%w: missing zone: %v", ErrInvalidHulls, zone)
		}
	}
	return nil
}

func (h Hulls) Values() iter.Seq[Hull] {
	return func(yield func(Hull) bool) {
		for _, hull := range *h.data {
			if !yield(hull) {
				break
			}
		}
	}
}

func (h Hulls) At(zone Zone) Hull {
	return (*h.data)[int(zone)-1]
}
