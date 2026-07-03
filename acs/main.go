// Command acs is a minimal TR-069 (CWMP) Auto Configuration Server.
//
// It accepts CPE Inform messages at POST /acs, stores the reporting device in
// SQLite, and exposes the collected devices as JSON for the Next.js device list UI.
package main

import (
	"log"
	"net/http"
	"os"

	"gtr069/acs/internal/httpapi"
	"gtr069/acs/internal/store"
)

func main() {
	addr := envOr("ACS_ADDR", ":7547")
	dbPath := envOr("ACS_DB", "acs.db")

	st, err := store.Open(dbPath)
	if err != nil {
		log.Fatalf("open store: %v", err)
	}
	defer st.Close()

	srv := &httpapi.Server{Store: st}
	mux := http.NewServeMux()
	srv.Register(mux)

	log.Printf("TR-069 ACS listening on %s (db=%s)", addr, dbPath)
	log.Printf("  Inform endpoint: POST %s/acs", addr)
	log.Printf("  Device API:      GET  %s/api/devices", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
