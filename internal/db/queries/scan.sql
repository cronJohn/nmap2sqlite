-- name: UpdateScanInfo :exec
UPDATE scans
SET scan_info_type = ?, scan_info_protocol = ?, scan_info_num_services = ?, scan_info_services = ?
WHERE id = ?;

-- name: UpdateVerboseInfo :exec
UPDATE scans
SET verbose_level = ?
WHERE id = ?;

-- name: UpdateDebuggingInfo :exec
UPDATE scans
SET debugging_level = ?
WHERE id = ?;

-- name: UpdateRunstatsInfo :exec
UPDATE scans
SET runstats_time = ?, runstats_time_str = ?, runstats_elapsed = ?, runstats_summary = ?, runstats_exit = ?, runstats_err = ?, 
hosts_up = ?, hosts_down = ?, hosts_total = ?
WHERE id = ?;
