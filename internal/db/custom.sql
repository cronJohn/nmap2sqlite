/* Custom sqlc queries that I will change the generated code for
because these specifically need to return the last insert id
These will just be used as a base in case I need to regen the code
*/

-- name: InsertNmaprunInfo :exec
INSERT INTO scans (
    scanner, args, start, start_str, version, xml_output_version
) VALUES (?, ?, ?, ?, ?, ?);

-- name: InsertHost :exec
INSERT INTO hosts (
    scan_id, start_time, end_time, status_state, status_reason, status_reason_ttl,
    distance_value, uptime_seconds, uptime_lastboot,
    tcp_sequence_index, tcp_sequence_difficulty, tcp_sequence_values,
    ip_id_seq_class, ip_id_seq_values,
    tcp_ts_seq_class, tcp_ts_seq_values,
    time_srtt, time_rttvar, time_to
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: InsertHostScript :exec
INSERT INTO host_scripts (
    host_id, script_id, script_output
) VALUES (?, ?, ?);

-- name: InsertHostTrace :exec
INSERT INTO host_traces (
    host_id, port, proto
) VALUES (?, ?, ?);

-- name: InsertPort :exec
INSERT INTO ports (
    host_id, protocol, port_id, 
    state, state_reason, state_reason_ttl, 
    service_name, service_conf, service_method, service_version, service_product, 
    service_extra_info, service_tunnel, service_proto, service_rpc_num, 
    service_lowver, service_highver, service_fp
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: InsertPortScript :exec
INSERT INTO port_scripts (
    port_id, script_id, script_output
) VALUES (?, ?, ?);

-- name: InsertExtraPorts :exec
INSERT INTO extra_ports (
    host_id, state, count
) VALUES (?, ?, ?);

-- name: InsertOsMatch :exec
INSERT INTO os_matches (
    host_id, name, accuracy, line
) VALUES (?, ?, ?, ?);
