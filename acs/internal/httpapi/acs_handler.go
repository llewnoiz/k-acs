// Package httpapi wires the ACS endpoint and the device REST API to the store.
package httpapi

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"gtr069/acs/internal/cwmp"
	"gtr069/acs/internal/store"
)

// Server holds shared dependencies for the HTTP handlers.
type Server struct {
	Store *store.Store
}

// Register attaches all routes to the given mux.
func (s *Server) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /acs", s.handleACS)
	mux.HandleFunc("GET /api/devices", withCORS(s.handleListDevices))
	mux.HandleFunc("GET /api/devices/{serial}", withCORS(s.handleGetDevice))
	mux.HandleFunc("OPTIONS /api/devices", withCORS(nil))
	mux.HandleFunc("OPTIONS /api/devices/{serial}", withCORS(nil))
}

// handleACS processes a single TR-069 CWMP request.
//
// A CPE session looks like: POST(Inform) -> InformResponse, then POST(empty) -> 204.
// This POC only reacts to the Inform; any empty POST ends the session.
func (s *Server) handleACS(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body", http.StatusBadRequest)
		return
	}

	// Empty POST => CPE has nothing more to send; close the session.
	if len(body) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	inform, cwmpID, err := cwmp.ParseInform(body)
	if err != nil {
		log.Printf("acs: failed to parse Inform: %v", err)
		http.Error(w, "invalid Inform", http.StatusInternalServerError)
		return
	}

	d := deviceFromInform(inform, clientIP(r))
	if err := s.Store.UpsertDevice(d); err != nil {
		log.Printf("acs: upsert failed: %v", err)
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}
	log.Printf("acs: Inform from %s (%s) event=%s ip=%s",
		d.SerialNumber, d.Manufacturer, d.LastEvent, d.IP)

	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Write(cwmp.BuildInformResponse(cwmpID))
}

// deviceFromInform maps an Inform onto a store.Device, pulling common parameters
// via suffix match so it works across "Device." and "InternetGatewayDevice." roots.
func deviceFromInform(in *cwmp.Inform, ip string) store.Device {
	if connURL := in.ParamSuffix("ManagementServer.ConnectionRequestURL"); connURL != "" {
		if host := hostFromURL(connURL); host != "" {
			ip = host
		}
	}

	params := map[string]string{}
	for _, p := range in.ParameterList {
		params[p.Name] = p.Value
	}
	paramJSON, _ := json.Marshal(params)

	now := time.Now().Format(time.RFC3339)
	return store.Device{
		SerialNumber:    in.DeviceID.SerialNumber,
		Manufacturer:    in.DeviceID.Manufacturer,
		OUI:             in.DeviceID.OUI,
		ProductClass:    in.DeviceID.ProductClass,
		IP:              ip,
		SoftwareVersion: in.ParamSuffix("DeviceInfo.SoftwareVersion"),
		HardwareVersion: in.ParamSuffix("DeviceInfo.HardwareVersion"),
		LastEvent:       in.FirstEventCode(),
		LastInformAt:    now,
		ParamJSON:       string(paramJSON),
		CreatedAt:       now,
	}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// hostFromURL extracts the host portion of a ConnectionRequestURL like
// "http://192.168.0.10:7547/cwmp". Returns "" if it cannot be parsed.
func hostFromURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return u.Hostname()
}
