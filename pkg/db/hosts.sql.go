// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: hosts.sql

package sqlc

import (
	"context"
)

const insertHostAddress = `-- name: InsertHostAddress :exec
INSERT INTO host_addresses (
    host_id, addr, addr_type, vendor
) VALUES (?, ?, ?, ?)
`

type InsertHostAddressParams struct {
	HostID   int64  `json:"host_id"`
	Addr     string `json:"addr"`
	AddrType string `json:"addr_type"`
	Vendor   string `json:"vendor"`
}

func (q *Queries) InsertHostAddress(ctx context.Context, arg InsertHostAddressParams) error {
	_, err := q.db.ExecContext(ctx, insertHostAddress,
		arg.HostID,
		arg.Addr,
		arg.AddrType,
		arg.Vendor,
	)
	return err
}

const insertHostNames = `-- name: InsertHostNames :exec
INSERT INTO host_names (
    host_id, name, type
) VALUES (?, ?, ?)
`

type InsertHostNamesParams struct {
	HostID int64  `json:"host_id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

func (q *Queries) InsertHostNames(ctx context.Context, arg InsertHostNamesParams) error {
	_, err := q.db.ExecContext(ctx, insertHostNames, arg.HostID, arg.Name, arg.Type)
	return err
}

const insertHostScriptElement = `-- name: InsertHostScriptElement :exec
INSERT INTO host_script_elements (
    host_script_id, key, value
) VALUES (?, ?, ?)
`

type InsertHostScriptElementParams struct {
	HostScriptID int64  `json:"host_script_id"`
	Key          string `json:"key"`
	Value        string `json:"value"`
}

func (q *Queries) InsertHostScriptElement(ctx context.Context, arg InsertHostScriptElementParams) error {
	_, err := q.db.ExecContext(ctx, insertHostScriptElement, arg.HostScriptID, arg.Key, arg.Value)
	return err
}

const insertTraceHop = `-- name: InsertTraceHop :exec
INSERT INTO trace_hops (
    trace_id, ttl, ip_addr, rtt, host
) VALUES (?, ?, ?, ?, ?)
`

type InsertTraceHopParams struct {
	TraceID int64   `json:"trace_id"`
	Ttl     int64   `json:"ttl"`
	IpAddr  string  `json:"ip_addr"`
	Rtt     float64 `json:"rtt"`
	Host    string  `json:"host"`
}

func (q *Queries) InsertTraceHop(ctx context.Context, arg InsertTraceHopParams) error {
	_, err := q.db.ExecContext(ctx, insertTraceHop,
		arg.TraceID,
		arg.Ttl,
		arg.IpAddr,
		arg.Rtt,
		arg.Host,
	)
	return err
}
