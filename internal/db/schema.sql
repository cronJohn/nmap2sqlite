CREATE TABLE IF NOT EXISTS scans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scanner TEXT NOT NULL,
    args TEXT NOT NULL,
    start INTEGER NOT NULL,
    start_str TEXT NOT NULL,
    version TEXT NOT NULL,
    xml_output_version TEXT NOT NULL,
    scan_info_type TEXT,
    scan_info_protocol TEXT,
    scan_info_num_services INTEGER,
    scan_info_services TEXT,
    verbose_level INTEGER,
    debugging_level INTEGER,
    runstats_time INTEGER,
    runstats_time_str TEXT,
    runstats_elapsed REAL,
    runstats_summary TEXT,
    runstats_exit TEXT,
    runstats_err TEXT,
    hosts_up INTEGER,
    hosts_down INTEGER,
    hosts_total INTEGER 
);

CREATE TABLE IF NOT EXISTS hosts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scan_id INTEGER NOT NULL,
    start_time INTEGER NOT NULL,
    end_time INTEGER NOT NULL,
    status_state TEXT NOT NULL,
    status_reason TEXT NOT NULL,
    status_reason_ttl INTEGER NOT NULL,
    distance_value INTEGER NOT NULL,
    uptime_seconds INTEGER NOT NULL,
    uptime_lastboot TEXT NOT NULL,
    tcp_sequence_index INTEGER NOT NULL,
    tcp_sequence_difficulty TEXT NOT NULL,
    tcp_sequence_values TEXT NOT NULL,
    ip_id_seq_class TEXT NOT NULL,
    ip_id_seq_values TEXT NOT NULL,
    tcp_ts_seq_class TEXT,
    tcp_ts_seq_values TEXT,
    time_srtt INTEGER NOT NULL,
    time_rttvar INTEGER NOT NULL,
    time_to INTEGER NOT NULL,
    FOREIGN KEY(scan_id) REFERENCES scans(id)
);

CREATE TABLE IF NOT EXISTS host_addresses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    host_id INTEGER NOT NULL,
    addr TEXT NOT NULL,
    addr_type TEXT NOT NULL,
    vendor TEXT NOT NULL,
    FOREIGN KEY(host_id) REFERENCES hosts(id)
);

CREATE TABLE IF NOT EXISTS host_names (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    host_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    FOREIGN KEY(host_id) REFERENCES hosts(id)
);

CREATE TABLE IF NOT EXISTS host_scripts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    host_id INTEGER NOT NULL,
    script_id TEXT NOT NULL,
    script_output TEXT NOT NULL,
    FOREIGN KEY(host_id) REFERENCES hosts(id)
);

CREATE TABLE IF NOT EXISTS host_script_elements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    host_script_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY(host_script_id) REFERENCES host_scripts(id)
);

CREATE TABLE IF NOT EXISTS host_traces (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    host_id INTEGER NOT NULL,
    port INTEGER NOT NULL,
    proto TEXT NOT NULL,
    FOREIGN KEY(host_id) REFERENCES hosts(id)
    );

CREATE TABLE IF NOT EXISTS trace_hops (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    trace_id INTEGER NOT NULL,
    ttl INTEGER NOT NULL,
    ip_addr TEXT NOT NULL,
    rtt REAL NOT NULL,
    host TEXT NOT NULL,
    FOREIGN KEY(trace_id) REFERENCES host_traces(id)
    );

CREATE TABLE IF NOT EXISTS ports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    host_id INTEGER NOT NULL,
    protocol TEXT NOT NULL,
    port_id INTEGER NOT NULL, -- External scanned port
    state TEXT NOT NULL,
    state_reason TEXT NOT NULL,
    state_reason_ttl INTEGER NOT NULL,
    service_name TEXT NOT NULL,
    service_conf INTEGER NOT NULL,
    service_method TEXT NOT NULL,
    service_version TEXT NOT NULL,
    service_product TEXT NOT NULL,
    service_extra_info TEXT NOT NULL,
    service_tunnel TEXT NOT NULL,
    service_proto TEXT NOT NULL,
    service_rpc_num TEXT NOT NULL,
    service_lowver TEXT NOT NULL,
    service_highver TEXT NOT NULL,
    service_fp TEXT NOT NULL,
    FOREIGN KEY(host_id) REFERENCES hosts(id)
);

CREATE TABLE IF NOT EXISTS port_service_cpes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    port_id INTEGER NOT NULL,
    cpe TEXT NOT NULL,
    FOREIGN KEY(port_id) REFERENCES ports(id)
);

CREATE TABLE IF NOT EXISTS port_scripts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    port_id INTEGER NOT NULL,
    script_id TEXT NOT NULL,
    script_output TEXT NOT NULL,
    FOREIGN KEY(port_id) REFERENCES ports(id)
);

CREATE TABLE IF NOT EXISTS port_script_elements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    port_script_id INTEGER NOT NULL,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    FOREIGN KEY(port_script_id) REFERENCES port_scripts(id)
);

CREATE TABLE IF NOT EXISTS extra_ports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    host_id INTEGER NOT NULL,
    state TEXT NOT NULL,
    count INTEGER NOT NULL,
    FOREIGN KEY(host_id) REFERENCES hosts(id)
);

CREATE TABLE IF NOT EXISTS extra_ports_reasons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    extra_port_id INTEGER NOT NULL,
    reason TEXT NOT NULL,
    count INTEGER NOT NULL,
    proto TEXT NOT NULL,
    ports TEXT NOT NULL,
    FOREIGN KEY(extra_port_id) REFERENCES extra_ports(id)
);

CREATE TABLE IF NOT EXISTS os_matches (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    host_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    accuracy TEXT NOT NULL,
    line TEXT NOT NULL,
    FOREIGN KEY(host_id) REFERENCES hosts(id)
);

CREATE TABLE IF NOT EXISTS os_classes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    os_match_id INTEGER NOT NULL,
    type TEXT NOT NULL,
    vendor TEXT NOT NULL,
    osfamily TEXT NOT NULL,
    osgen TEXT NOT NULL,
    accuracy TEXT NOT NULL,
    cpe TEXT NOT NULL,
    FOREIGN KEY(os_match_id) REFERENCES os_matches(id)
);

CREATE TABLE IF NOT EXISTS os_ports_used (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    host_id INTEGER NOT NULL,
    state TEXT NOT NULL,
    protocol TEXT NOT NULL,
    port_id INTEGER NOT NULL, -- Internal fingerprinting hint port
    FOREIGN KEY(host_id) REFERENCES hosts(id)
);
