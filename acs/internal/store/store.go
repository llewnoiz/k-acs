// Package store persists CPE devices seen via Inform into an embedded SQLite database.
package store

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "modernc.org/sqlite" // pure-Go SQLite driver (no CGO)
)

//go:embed schema.sql
var schema string

// Device is a CPE as tracked by the ACS.
type Device struct {
	SerialNumber    string `json:"serialNumber"`
	Manufacturer    string `json:"manufacturer"`
	OUI             string `json:"oui"`
	ProductClass    string `json:"productClass"`
	IP              string `json:"ip"`
	SoftwareVersion string `json:"softwareVersion"`
	HardwareVersion string `json:"hardwareVersion"`
	LastEvent       string `json:"lastEvent"`
	LastInformAt    string `json:"lastInformAt"`
	InformCount     int    `json:"informCount"`
	ParamJSON       string `json:"paramJson,omitempty"`
	CreatedAt       string `json:"createdAt"`
}

// Store wraps the SQLite connection.
type Store struct {
	db *sql.DB
}

// Open opens (creating if needed) the SQLite database at path and applies the schema.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("apply schema: %w", err)
	}
	return &Store{db: db}, nil
}

// Close closes the underlying database.
func (s *Store) Close() error { return s.db.Close() }

// UpsertDevice inserts a new device or updates an existing one (keyed by serial number).
// created_at is preserved and inform_count is incremented on conflict.
func (s *Store) UpsertDevice(d Device) error {
	const q = `
INSERT INTO devices (
    serial_number, manufacturer, oui, product_class, ip,
    software_version, hardware_version, last_event, last_inform_at,
    inform_count, param_json, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1, ?, ?)
ON CONFLICT(serial_number) DO UPDATE SET
    manufacturer     = excluded.manufacturer,
    oui              = excluded.oui,
    product_class    = excluded.product_class,
    ip               = excluded.ip,
    software_version = excluded.software_version,
    hardware_version = excluded.hardware_version,
    last_event       = excluded.last_event,
    last_inform_at   = excluded.last_inform_at,
    inform_count     = devices.inform_count + 1,
    param_json       = excluded.param_json;`
	_, err := s.db.Exec(q,
		d.SerialNumber, d.Manufacturer, d.OUI, d.ProductClass, d.IP,
		d.SoftwareVersion, d.HardwareVersion, d.LastEvent, d.LastInformAt,
		d.ParamJSON, d.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert device %q: %w", d.SerialNumber, err)
	}
	return nil
}

const selectCols = `serial_number, manufacturer, oui, product_class, ip,
    software_version, hardware_version, last_event, last_inform_at,
    inform_count, param_json, created_at`

func scanDevice(sc interface{ Scan(...any) error }, d *Device) error {
	return sc.Scan(
		&d.SerialNumber, &d.Manufacturer, &d.OUI, &d.ProductClass, &d.IP,
		&d.SoftwareVersion, &d.HardwareVersion, &d.LastEvent, &d.LastInformAt,
		&d.InformCount, &d.ParamJSON, &d.CreatedAt,
	)
}

// ListDevices returns all devices, most recently seen first. ParamJSON is omitted from the list view.
func (s *Store) ListDevices() ([]Device, error) {
	rows, err := s.db.Query(`SELECT ` + selectCols + ` FROM devices ORDER BY last_inform_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}
	defer rows.Close()

	devices := []Device{}
	for rows.Next() {
		var d Device
		if err := scanDevice(rows, &d); err != nil {
			return nil, fmt.Errorf("scan device: %w", err)
		}
		d.ParamJSON = "" // keep the list payload small; full params via GetDevice
		devices = append(devices, d)
	}
	return devices, rows.Err()
}

// GetDevice returns a single device (including ParamJSON) or (nil, nil) if not found.
func (s *Store) GetDevice(serial string) (*Device, error) {
	row := s.db.QueryRow(`SELECT `+selectCols+` FROM devices WHERE serial_number = ?`, serial)
	var d Device
	if err := scanDevice(row, &d); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get device %q: %w", serial, err)
	}
	return &d, nil
}
