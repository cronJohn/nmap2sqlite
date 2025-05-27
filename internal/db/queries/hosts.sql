-- name: InsertHostNames :exec
INSERT INTO host_names (
    host_id, name, type
) VALUES (?, ?, ?);

-- name: InsertHostAddress :exec
INSERT INTO host_addresses (
    host_id, addr, addr_type, vendor
) VALUES (?, ?, ?, ?);

-- name: InsertHostScriptElement :exec
INSERT INTO host_script_elements (
    host_script_id, key, value
) VALUES (?, ?, ?);

-- name: InsertTraceHop :exec
INSERT INTO trace_hops (
    trace_id, ttl, ip_addr, rtt, host
) VALUES (?, ?, ?, ?, ?);

