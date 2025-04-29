package api

func (server *Server) drawRoutes() {
	server.router.HandleFunc(
		"GET /devices",
		server.GetDevices,
	)

	server.router.HandleFunc(
		"GET /devices/{device_id}",
		server.GetDevice,
	)

	server.router.HandleFunc(
		"GET /devices/{device_id}/services/{service_uuid}/characteristics/{char_uuid}",
		server.ReadCharacteristic,
	)

	server.router.HandleFunc(
		"POST /devices/{device_id}/services/{service_uuid}/characteristics/{char_uuid}",
		server.WriteCharacteristic,
	)

	server.router.HandleFunc(
		"POST /start",
		server.StartScan,
	)

	server.router.HandleFunc(
		"POST /stop",
		server.StopScan,
	)
}
