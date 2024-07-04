package filters

import (
	"regexp"

	"tinygo.org/x/bluetooth"
)

type NameFilter struct {
	r *regexp.Regexp
}

func (filter *NameFilter) Filter(device bluetooth.ScanResult) bool {
	return filter.r.MatchString(device.LocalName())
}
