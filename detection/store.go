package detection

type Store struct {
	devices map[string]Device
}

func NewStore() *Store {
	return &Store{
		devices: map[string]Device{},
	}
}

func (store *Store) Watch(detectionChannel chan Device, quit chan bool) error {
	for {
		select {
		case device := <-detectionChannel:
			store.Save(device)
		case <-quit:
			return nil
		}
	}
}

func (store *Store) Save(device Device) {
	store.devices[device.Address.String()] = device
}

func (store *Store) Get(address string) (Device, bool) {
	device, ok := store.devices[address]
	return device, ok
}

func (store *Store) List() []Device {
	var devices []Device

	for _, device := range store.devices {
		devices = append(devices, device)
	}

	return devices
}
