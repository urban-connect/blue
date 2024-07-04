package discovery

import (
	"fmt"

	"tinygo.org/x/bluetooth"
)

type Connection struct {
	Name    string
	Address bluetooth.Address

	bleDevice bluetooth.Device
}

func (connection *Connection) Read(serviceUUID bluetooth.UUID, charUUID bluetooth.UUID) ([]byte, error) {
	services, err := connection.bleDevice.DiscoverServices([]bluetooth.UUID{
		serviceUUID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find service: %w", err)
	}

	chars, err := services[0].DiscoverCharacteristics([]bluetooth.UUID{
		charUUID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to find characteristic: %w", err)
	}

	data := make([]byte, 32)

	_, err = chars[0].Read(data)

	if err != nil {
		return nil, fmt.Errorf("failed to read characteristic: %w", err)
	}

	return data, nil
}

func (connection *Connection) Discover() (DiscoveredServices, error) {
	services, err := connection.bleDevice.DiscoverServices(nil)

	if err != nil {
		return nil, fmt.Errorf("failed to discover services: %w", err)
	}

	discovery := DiscoveredServices{}

	for i := range services {
		chars, err := services[i].DiscoverCharacteristics(nil)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to discover characteristics for %s: %w",
				services[i].UUID().String(),
				err,
			)
		}

		var charUUIDs []string

		for j := range chars {
			charUUIDs = append(charUUIDs, chars[j].UUID().String())
		}

		discovery[services[i].UUID().String()] = charUUIDs
	}

	return discovery, nil
}

func (connection *Connection) Disconnect() error {
	return connection.bleDevice.Disconnect()
}
