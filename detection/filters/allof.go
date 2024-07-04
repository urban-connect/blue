package filters

import "tinygo.org/x/bluetooth"

type AllOfFilter struct {
	filters []Filter
}

func (filter *AllOfFilter) Filter(device bluetooth.ScanResult) bool {
	for _, f := range filter.filters {
		if !f.Filter(device) {
			return false
		}
	}

	return true
}
