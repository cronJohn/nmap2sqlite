-- name: InsertOsClass :exec
INSERT INTO os_classes (
    os_match_id, type, vendor, osfamily, osgen, accuracy, cpe
) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: InsertOsPortsUsed :exec
INSERT INTO os_ports_used (
    host_id, state, protocol, port_id
) VALUES (?, ?, ?, ?);
