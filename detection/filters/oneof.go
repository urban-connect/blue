package filters

import "tinygo.org/x/bluetooth"

type OneOfFilter struct {
	filters []Filter
}

func (filter *OneOfFilter) Filter(device bluetooth.ScanResult) bool {
	for _, f := range filter.filters {
		if f.Filter(device) {
			return true
		}
	}

	return false
}
