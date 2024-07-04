package discovery

import (
	"fmt"
	"time"

	"tinygo.org/x/bluetooth"
	"urban-connect.ch/blue/detection"
)

type DiscoveredServices map[string][]string

func Connect(device detection.Device) (*Connection, error) {
	bleDevice, err := bluetooth.DefaultAdapter.Connect(
		device.Address,
		bluetooth.ConnectionParams{
			ConnectionTimeout: bluetooth.NewDuration(60 * time.Second),
			Timeout:           bluetooth.NewDuration(60 * time.Second),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", device.Address.String(), err)
	}

	connection := &Connection{
		Name:      device.Name,
		Address:   device.Address,
		bleDevice: bleDevice,
	}

	return connection, nil
}
