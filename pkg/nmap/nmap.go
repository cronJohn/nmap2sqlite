package nmap

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/rs/zerolog/log"

	sqlc "github.com/cronJohn/nmap2sqlite/pkg/db"
	_ "github.com/cronJohn/nmap2sqlite/pkg/logger"
)

func ParseData(ctx context.Context, data io.Reader, dbHandle *sql.DB) error {
	if dbHandle == nil {
		return errors.New("No database handle provided")
	}
	var (
		tx     *sql.Tx
		q      *sqlc.Queries
		scanID int64 // Store scanID to link scan information with the corresponding scan record

	)
	decoder := xml.NewDecoder(data)

	for {
		select {
		case <-ctx.Done(): // If cancelled early
			return ctx.Err()
		default:
			token, err := decoder.Token()
			if err == io.EOF {
				log.Debug().Msg("No more xml to parse")
				return nil
			}
			if err != nil {
				log.Error().Err(err).Msg("Error decoding token")
				continue
			}

			if startElem, ok := token.(xml.StartElement); ok {
				log.Debug().Msgf("Received element: <%s>", startElem.Name.Local)
				switch startElem.Name.Local {
				case "nmaprun":
					// Use map instead of raw indexing for more expressive code
					attrMap := make(map[string]string, 6)

					for _, attr := range startElem.Attr {
						attrMap[attr.Name.Local] = attr.Value
					}

					// Start transaction for init data
					tx, err = dbHandle.BeginTx(ctx, nil)
					if err != nil {
						log.Error().Err(err).Msg("Error starting transaction for init data")
						return err
					}
					defer func() {
						if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
							log.Error().Err(err).Msg("Rollback failed")
						}
					}()

					q = sqlc.New(tx)
					// Extract values safely using the map
					lastReturnId, err := q.InsertNmaprunInfo(ctx, sqlc.InsertNmaprunInfoParams{
						Scanner:          attrMap["scanner"],
						Args:             attrMap["args"],
						Start:            convStringToInt(attrMap["start"]),
						StartStr:         attrMap["startstr"],
						Version:          attrMap["version"],
						XmlOutputVersion: attrMap["xmloutputversion"],
					})
					if err != nil {
						log.Error().Err(err).Msg("Error inserting nmaprun info")
					}

					scanID = lastReturnId
				case "scaninfo":
					scanInfo := ScanInfo{}
					if err := decoder.DecodeElement(&scanInfo, &startElem); err != nil {
						log.Error().Err(err).Msg("Error decoding scaninfo")
						continue
					}

					q.UpdateScanInfo(ctx, sqlc.UpdateScanInfoParams{
						ScanInfoType:     sql.NullString{String: scanInfo.Type, Valid: true},
						ScanInfoProtocol: sql.NullString{String: scanInfo.Protocol, Valid: true},
						ScanInfoNumServices: sql.NullInt64{
							Int64: scanInfo.NumServices,
							Valid: true,
						},
						ScanInfoServices: sql.NullString{String: scanInfo.Services, Valid: true},
						ID:               scanID,
					})

				case "verbose":
					verboseInfo := Verbose{}
					if err := decoder.DecodeElement(&verboseInfo, &startElem); err != nil {
						log.Error().Err(err).Msg("Error decoding verbose info")
						continue
					}

					q.UpdateVerboseInfo(ctx, sqlc.UpdateVerboseInfoParams{
						VerboseLevel: sql.NullInt64{Int64: verboseInfo.Level, Valid: true},
						ID:           scanID,
					})
				case "debugging":
					debuggingInfo := Debugging{}
					if err := decoder.DecodeElement(&debuggingInfo, &startElem); err != nil {
						log.Error().Err(err).Msg("Error decoding debugging info")
						continue
					}

					q.UpdateDebuggingInfo(ctx, sqlc.UpdateDebuggingInfoParams{
						DebuggingLevel: sql.NullInt64{Int64: debuggingInfo.Level, Valid: true},
						ID:             scanID,
					})

					// Commit initial transaction data
					if err := tx.Commit(); err != nil {
						log.Error().Err(err).Msg("Error committing init transaction data")
						return err
					}
				case "host":
					host := Host{}
					if err := decoder.DecodeElement(&host, &startElem); err != nil {
						log.Error().Err(err).Msg("Error decoding host info")
						continue
					}

					if err := insertHostInfoIntoDB(ctx, host, scanID, dbHandle); err != nil {
						log.Error().Err(err).Msg("Error inserting host info")
						return err
					}
				case "runstats":
					runStats := RunStats{}
					if err := decoder.DecodeElement(&runStats, &startElem); err != nil {
						log.Error().Err(err).Msg("Error decoding runstats info")
						continue
					}

					// Refersh q since old tx transaction is closed by now
					q = sqlc.New(dbHandle)

					q.UpdateRunstatsInfo(ctx, sqlc.UpdateRunstatsInfoParams{
						RunstatsTime: sql.NullInt64{
							Int64: int64(runStats.Finished.Time),
							Valid: true,
						},
						RunstatsTimeStr: sql.NullString{
							String: runStats.Finished.TimeStr,
							Valid:  true,
						},
						RunstatsElapsed: sql.NullFloat64{
							Float64: runStats.Finished.Elapsed,
							Valid:   true,
						},
						RunstatsSummary: sql.NullString{
							String: runStats.Finished.Summary,
							Valid:  true,
						},
						RunstatsExit: sql.NullString{
							String: runStats.Finished.Exit,
							Valid:  true,
						},
						RunstatsErr: sql.NullString{
							String: runStats.Finished.ErrorMsg,
							Valid:  true,
						},
						HostsUp:    sql.NullInt64{Int64: runStats.Hosts.Up, Valid: true},
						HostsDown:  sql.NullInt64{Int64: runStats.Hosts.Down, Valid: true},
						HostsTotal: sql.NullInt64{Int64: runStats.Hosts.Total, Valid: true},
						ID:         scanID,
					})
				}
			}
		}
	}
}

// For simple one liners
func convStringToInt(s string) int64 {
	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Error converting string to time")
	}
	return t
}

func insertHostInfoIntoDB(
	ctx context.Context,
	host Host,
	scanID int64,
	dbHandle *sql.DB,
) error {
	tx, err := dbHandle.BeginTx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction for host info")
		return err
	}
	defer tx.Rollback()

	q := sqlc.New(tx)

	// Insert host info
	hostID, err := q.InsertHost(ctx, sqlc.InsertHostParams{
		ScanID:                scanID,
		StartTime:             int64(host.StartTime),
		EndTime:               int64(host.EndTime),
		StatusState:           host.Status.State,
		StatusReason:          host.Status.Reason,
		StatusReasonTtl:       host.Status.ReasonTTL,
		DistanceValue:         host.Distance.Value,
		UptimeSeconds:         host.Uptime.Seconds,
		UptimeLastboot:        host.Uptime.Lastboot,
		TcpSequenceIndex:      host.TcpSequence.Index,
		TcpSequenceDifficulty: host.TcpSequence.Difficulty,
		TcpSequenceValues:     host.TcpSequence.Values,
		IpIDSeqClass:          host.IpIdSequence.Class,
		IpIDSeqValues:         host.IpIdSequence.Values,
		TcpTsSeqClass:         host.TcpTsSequence.Class,
		TcpTsSeqValues:        host.TcpTsSequence.Values,
		TimeSrtt:              host.Times.SRTT,
		TimeRttvar:            host.Times.RTTVAR,
		TimeTo:                host.Times.To,
	})
	if err != nil {
		return err
	}

	for _, hostname := range host.Hostnames {
		q.InsertHostNames(ctx, sqlc.InsertHostNamesParams{
			HostID: hostID,
			Name:   hostname.Name,
			Type:   hostname.Type,
		})
	}

	for _, address := range host.Addresses {
		q.InsertHostAddress(ctx, sqlc.InsertHostAddressParams{
			HostID:   hostID,
			Addr:     address.Addr,
			AddrType: address.AddrType,
			Vendor:   address.Vendor,
		})
	}

	for _, script := range host.HostScripts {
		scriptID, err := q.InsertHostScript(ctx, sqlc.InsertHostScriptParams{
			HostID:       hostID,
			ScriptID:     script.Id,
			ScriptOutput: script.Output,
		})
		if err != nil {
			return err
		}

		for _, element := range script.Elements {
			err = q.InsertHostScriptElement(ctx, sqlc.InsertHostScriptElementParams{
				HostScriptID: scriptID,
				Key:          element.Key,
				Value:        element.Value,
			})
			if err != nil {
				return err
			}
		}
	}

	hostTracesID, err := q.InsertHostTrace(ctx, sqlc.InsertHostTraceParams{
		HostID: hostID,
		Port:   host.Trace.Port,
		Proto:  host.Trace.Proto,
	})

	for _, hop := range host.Trace.Hops {
		q.InsertTraceHop(ctx, sqlc.InsertTraceHopParams{
			TraceID: hostTracesID,
			Ttl:     hop.TTL,
			IpAddr:  hop.IPAddr,
			Rtt:     hop.RTT,
			Host:    hop.Host,
		})
	}

	// Insert port info
	for _, port := range host.Ports {
		portID, err := q.InsertPort(ctx, sqlc.InsertPortParams{
			HostID:           hostID,
			Protocol:         port.Protocol,
			PortID:           port.PortId,
			State:            port.State.State,
			StateReason:      port.State.Reason,
			StateReasonTtl:   port.State.ReasonTTL,
			ServiceName:      port.Service.Name,
			ServiceConf:      port.Service.Conf,
			ServiceMethod:    port.Service.Method,
			ServiceVersion:   port.Service.Version,
			ServiceProduct:   port.Service.Product,
			ServiceExtraInfo: port.Service.ExtraInfo,
			ServiceTunnel:    port.Service.Tunnel,
			ServiceProto:     port.Service.Proto,
			ServiceRpcNum:    port.Service.Rpcnum,
			ServiceLowver:    port.Service.Lowver,
			ServiceHighver:   port.Service.Highver,
			ServiceFp:        port.Service.ServiceFp,
		})
		if err != nil {
			return err
		}

		q.InsertPortServiceCpe(ctx, sqlc.InsertPortServiceCpeParams{
			PortID: portID,
			Cpe:    fmt.Sprintf("%s", port.Service.CPEs),
		})

		for _, portScript := range port.Scripts {
			portScriptID, err := q.InsertPortScript(ctx, sqlc.InsertPortScriptParams{
				PortID:       portID,
				ScriptID:     portScript.Id,
				ScriptOutput: portScript.Output,
			})
			if err != nil {
				return err
			}

			for _, element := range portScript.Elements {
				err = q.InsertPortScriptElement(ctx, sqlc.InsertPortScriptElementParams{
					PortScriptID: portScriptID,
					Key:          element.Key,
					Value:        element.Value,
				})
				if err != nil {
					return err
				}
			}
		}
	}

	// Insert extra ports
	for _, extraPort := range host.ExtraPorts {
		extraPortID, err := q.InsertExtraPorts(ctx, sqlc.InsertExtraPortsParams{
			HostID: hostID,
			State:  extraPort.State,
			Count:  extraPort.Count,
		})
		if err != nil {
			return err
		}

		for _, extraPortReason := range extraPort.ExtraReasons {
			err = q.InsertExtraPortsReason(ctx, sqlc.InsertExtraPortsReasonParams{
				ExtraPortID: extraPortID,
				Reason:      extraPortReason.Reason,
				Count:       extraPortReason.Count,
				Proto:       extraPortReason.Proto,
				Ports:       extraPortReason.Ports,
			})
			if err != nil {
				return err
			}
		}
	}

	// Insert os info
	for _, osMatch := range host.Os.OsMatches {
		osMatchID, err := q.InsertOsMatch(ctx, sqlc.InsertOsMatchParams{
			HostID:   hostID,
			Name:     osMatch.Name,
			Accuracy: osMatch.Accuracy,
			Line:     osMatch.Line,
		})
		if err != nil {
			return err
		}
		for _, osClass := range osMatch.OsClasses {
			err = q.InsertOsClass(ctx, sqlc.InsertOsClassParams{
				OsMatchID: osMatchID,
				Type:      osClass.Type,
				Vendor:    osClass.Vendor,
				Osfamily:  osClass.OsFamily,
				Osgen:     osClass.OsGen,
				Accuracy:  osClass.Accuracy,
				Cpe:       fmt.Sprintf("%s", osClass.CPEs),
			})
			if err != nil {
				return err
			}
		}
	}

	for _, osPortsUsed := range host.Os.PortsUsed {
		q.InsertOsPortsUsed(ctx, sqlc.InsertOsPortsUsedParams{
			HostID:   hostID,
			State:    osPortsUsed.State,
			Protocol: osPortsUsed.Proto,
			PortID:   osPortsUsed.PortId,
		})
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("Error committing host info")
		return err
	}

	return nil
}
