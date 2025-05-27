-- name: InsertPortServiceCpe :exec
INSERT INTO port_service_cpes (
    port_id, cpe
) VALUES (?, ?);

-- name: InsertPortScriptElement :exec
INSERT INTO port_script_elements (
    port_script_id, key, value
) VALUES (?, ?, ?);

-- name: InsertExtraPortsReason :exec
INSERT INTO extra_ports_reasons (
    extra_port_id, reason, count, proto, ports
) VALUES (?, ?, ?, ?, ?);
