package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"tinygo.org/x/bluetooth"
	"urban-connect.ch/blue/detection"
	"urban-connect.ch/blue/discovery"
)

type Server struct {
	router           *http.ServeMux
	store            *detection.Store
	detectionChannel chan detection.Device
}

type Error struct {
	Message string `json:"message"`
}

func NewServer(store *detection.Store, detectionChannel chan detection.Device) *Server {
	server := &Server{
		router:           http.NewServeMux(),
		store:            store,
		detectionChannel: detectionChannel,
	}

	server.drawRoutes()

	return server
}

func (server *Server) StartScan(w http.ResponseWriter, r *http.Request) {
	var request StartScanRequest
	var builder Builder

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		server.RespondWithError(
			w,
			http.StatusUnprocessableEntity,
			fmt.Errorf("failed to parse request body: %w", err),
		)

		return
	}

	filter, err := builder.Build(request.Filter)

	if err != nil {
		server.RespondWithError(
			w,
			http.StatusUnprocessableEntity,
			fmt.Errorf("failed to build a filter from request body: %w", err),
		)

		return
	}

	go func() {
		scanner := detection.NewScanner(filter)

		fmt.Println("Scanner has been started...")

		if err := scanner.Run(server.detectionChannel); err != nil {
			fmt.Printf("Failed to start the scanner: %v\n", err)
		}

		fmt.Println("Scanner has been stopped...")
	}()

	server.RespondWith(w, http.StatusOK, map[string]string{
		"scanning": "started",
	})
}

func (server *Server) StopScan(w http.ResponseWriter, r *http.Request) {
	err := bluetooth.DefaultAdapter.StopScan()

	if err != nil {
		server.RespondWithError(
			w,
			http.StatusNotFound,
			fmt.Errorf("failed to stop scanning: %w", err),
		)

		return
	}

	server.RespondWith(w, http.StatusOK, map[string]string{
		"scanning": "stopped",
	})
}

func (server *Server) GetDevices(w http.ResponseWriter, r *http.Request) {
	devices := server.store.List()

	result := DeviceListResponse{
		Devices: make([]DeviceResponse, len(devices)),
	}

	for i, d := range devices {
		result.Devices[i] = DeviceResponse{
			Address: d.Address.String(),
			Name:    d.Name,
		}
	}

	server.RespondWith(w, http.StatusOK, result)
}

func (server *Server) GetDevice(w http.ResponseWriter, r *http.Request) {
	bleDevice, found := server.store.Get(r.PathValue("device_id"))

	if !found {
		server.RespondWithError(
			w,
			http.StatusNotFound,
			fmt.Errorf("device not found: %s", r.PathValue("device_id")),
		)

		return
	}

	connection, err := discovery.Connect(bleDevice)

	if err != nil {
		server.RespondWithError(
			w,
			http.StatusNotFound,
			fmt.Errorf("failed to connect: %w", err),
		)

		return
	}

	defer connection.Disconnect()

	discoveredServices, err := connection.Discover()

	if err != nil {
		server.RespondWithError(
			w,
			http.StatusInternalServerError,
			fmt.Errorf("failed to read data: %w", err),
		)

		return
	}

	server.RespondWith(w, http.StatusOK, map[string]interface{}{
		"discovered_services": discoveredServices,
	})
}

func (server *Server) ReadCharacteristic(w http.ResponseWriter, r *http.Request) {
	bleDevice, found := server.store.Get(r.PathValue("device_id"))

	if !found {
		server.RespondWithError(
			w,
			http.StatusNotFound,
			fmt.Errorf("device not found: %s", r.PathValue("device_id")),
		)

		return
	}

	serviceUUID, err := bluetooth.ParseUUID(r.PathValue("service_uuid"))

	if err != nil {
		server.RespondWithError(
			w,
			http.StatusUnprocessableEntity,
			fmt.Errorf("invalid service uuid: %s", r.PathValue("service_uuid")),
		)

		return
	}

	charUUID, err := bluetooth.ParseUUID(r.PathValue("char_uuid"))

	if err != nil {
		server.RespondWithError(
			w,
			http.StatusUnprocessableEntity,
			fmt.Errorf("invalid characteristic uuid: %s", r.PathValue("char_uuid")),
		)

		return
	}

	connection, err := discovery.Connect(bleDevice)

	if err != nil {
		server.RespondWithError(
			w,
			http.StatusNotFound,
			fmt.Errorf("failed to connect: %w", err),
		)

		return
	}

	defer connection.Disconnect()

	data, err := connection.Read(serviceUUID, charUUID)

	if err != nil {
		server.RespondWithError(
			w,
			http.StatusNotFound,
			fmt.Errorf("failed to read data: %w", err),
		)

		return
	}

	var result [32]string

	for i, d := range data {
		result[i] = fmt.Sprintf("0x%x", d)
	}

	server.RespondWith(w, http.StatusOK, map[string]interface{}{
		"data": result,
	})
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.router.ServeHTTP(w, r)
}

func (server *Server) RespondWithError(w http.ResponseWriter, status int, err error) {
	server.RespondWith(w, status, Error{
		Message: err.Error(),
	})
}

func (server *Server) RespondWith(w http.ResponseWriter, status int, data interface{}) {
	payload, err := json.Marshal(data)

	if err != nil {
		payload = nil
		status = http.StatusInternalServerError
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload != nil {
		w.Write(payload)
	}
}
