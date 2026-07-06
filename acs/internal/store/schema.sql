CREATE TABLE IF NOT EXISTS devices (
    serial_number    TEXT PRIMARY KEY,
    manufacturer     TEXT NOT NULL DEFAULT '',
    oui              TEXT NOT NULL DEFAULT '',
    product_class    TEXT NOT NULL DEFAULT '',
    ip               TEXT NOT NULL DEFAULT '',
    software_version TEXT NOT NULL DEFAULT '',
    hardware_version TEXT NOT NULL DEFAULT '',
    last_event       TEXT NOT NULL DEFAULT '',
    last_inform_at   TEXT NOT NULL DEFAULT '',
    inform_count     INTEGER NOT NULL DEFAULT 0,
    param_json       TEXT NOT NULL DEFAULT '{}',
    created_at       TEXT NOT NULL DEFAULT '',
    up_time          INTEGER NOT NULL DEFAULT 0
);
