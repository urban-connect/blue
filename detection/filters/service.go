package filters

import "tinygo.org/x/bluetooth"

type ServiceFilter struct {
	UUID bluetooth.UUID
}

func (filter *ServiceFilter) Filter(device bluetooth.ScanResult) bool {
	return device.HasServiceUUID(filter.UUID)
}
