package detection

import (
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
	"urban-connect.ch/blue/detection/filters"
)

type onDeviceDiscovered func(adapter *bluetooth.Adapter, bld bluetooth.ScanResult)

type Scanner struct {
	filter filters.Filter
}

func NewScanner(filter filters.Filter) *Scanner {
	return &Scanner{
		filter: filter,
	}
}

func (scanner *Scanner) Run(detectionChannel chan Device) error {
	callback := scanner.onDeviceDiscovered(detectionChannel)

	if err := bluetooth.DefaultAdapter.Scan(callback); err != nil {
		return fmt.Errorf("failed to start bluetooth scanner: %w", err)
	}

	return nil
}

func (scanner *Scanner) onDeviceDiscovered(detectionChannel chan Device) onDeviceDiscovered {
	return func(adapter *bluetooth.Adapter, bleDevice bluetooth.ScanResult) {
		if !scanner.filter.Filter(bleDevice) {
			return
		}

		fmt.Printf("Local name: %s\n", bleDevice.LocalName())
		fmt.Printf("Address: %s\n", bleDevice.Address.UUID.String())
		fmt.Println("----")

		detectedDevice := Device{
			Name:       bleDevice.LocalName(),
			Address:    bleDevice.Address,
			RSSI:       bleDevice.RSSI,
			DetectedAt: time.Now(),
		}

		detectionChannel <- detectedDevice
	}
}
