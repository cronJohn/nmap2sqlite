package nmap

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"os/exec"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

var (
	dbPath   = flag.String("d", "../../internal/db/test.sqlite", "Path to sqlite database")
	xmlPath  = flag.String("x", "./test.xml", "Path to nmap xml test file")
	nmapArgs = flag.String("n", "-sS localhost", "Nmap custom arguments")
	xmlFile  *os.File
	dbHandle *sql.DB
)

func TestMain(m *testing.M) {
	log.Info().Msgf("XML path: %v", *xmlPath)
	log.Info().Msgf("DB path: %v", *dbPath)

	var err error
	xmlFile, err = os.Open(*xmlPath)
	if err != nil {
		log.Error().Err(err).Msg("Error opening xml file")
		os.Exit(1)
	}

	defer xmlFile.Close()

	dbHandle, err = sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Error().Err(err).Msg("Error opening database")
		os.Exit(1)
	}
	defer dbHandle.Close()

	os.Exit(m.Run())
}

func TestXMLParseData(t *testing.T) {
	t.Run("Test XML parse normally", func(t *testing.T) {
		err := ParseData(context.Background(), xmlFile, dbHandle)
		if err != nil {
			t.Errorf("Error parsing data: %s", err)
		}
	})

	ctx, cancel := context.WithCancel(context.Background())

	t.Run("Test XML parse cancelling early", func(t *testing.T) {
		cancel()
		err := ParseData(ctx, xmlFile, dbHandle)
		if err != context.Canceled {
			t.Errorf("Error cancelling function: %s", err)
		}
	})
}

func TestNmapScanParseData(t *testing.T) {
	args := make([]string, 0, 5)
	args = []string{"-oX", "-"}
	args = append(args, strings.Fields(*nmapArgs)...)
	cmd := exec.Command("nmap", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Errorf("Error getting stdout pipe for nmap command: %s", err)
	}

	if err := cmd.Start(); err != nil {
		t.Errorf("Error starting nmap: %s", err)
	}

	err = ParseData(context.Background(), stdout, dbHandle)
	if err != nil {
		t.Errorf("Error parsing data: %s", err)
	}

	if err := cmd.Wait(); err != nil {
		t.Errorf("Error waiting for nmap: %s", err)
	}
}

func BenchmarkXMLParseData(b *testing.B) {
	ctx := context.Background()

	for b.Loop() {
		err := ParseData(ctx, xmlFile, dbHandle)
		if err != nil {
			b.Errorf("Error parsing data: %s", err)
		}
		xmlFile.Seek(0, 0)
	}
}
