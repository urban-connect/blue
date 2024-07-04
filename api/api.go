package api

type DeviceResponse struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type DeviceListResponse struct {
	Devices []DeviceResponse `json:"devices"`
}

type StartScanRequest struct {
	Filter Filter `json:"filter"`
}
