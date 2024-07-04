package filters

import (
	"regexp"

	"tinygo.org/x/bluetooth"
)

type Filter interface {
	Filter(device bluetooth.ScanResult) bool
}

func Name(r *regexp.Regexp) *NameFilter {
	return &NameFilter{
		r: r,
	}
}

func OneOf(filters ...Filter) *OneOfFilter {
	return &OneOfFilter{
		filters: filters,
	}
}

func AllOf(filters ...Filter) *AllOfFilter {
	return &AllOfFilter{
		filters: filters,
	}
}

func Service(uuid bluetooth.UUID) *ServiceFilter {
	return &ServiceFilter{
		UUID: uuid,
	}
}
