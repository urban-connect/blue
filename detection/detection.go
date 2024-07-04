package detection

import (
	"time"

	"tinygo.org/x/bluetooth"
)

type Device struct {
	Name       string
	Address    bluetooth.Address
	RSSI       int16
	DetectedAt time.Time
}
