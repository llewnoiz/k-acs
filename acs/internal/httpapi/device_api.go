package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("api: encode response: %v", err)
	}
}

// handleListDevices returns all known devices as a JSON array.
func (s *Server) handleListDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := s.Store.ListDevices()
	if err != nil {
		log.Printf("api: list devices: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "list failed"})
		return
	}
	writeJSON(w, http.StatusOK, devices)
}

// handleGetDevice returns a single device (with its full parameter map) by serial number.
func (s *Server) handleGetDevice(w http.ResponseWriter, r *http.Request) {
	serial := r.PathValue("serial")
	device, err := s.Store.GetDevice(serial)
	if err != nil {
		log.Printf("api: get device: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "get failed"})
		return
	}
	if device == nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "device not found"})
		return
	}
	writeJSON(w, http.StatusOK, device)
}
