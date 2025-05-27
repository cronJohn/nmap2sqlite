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
)

var (
	dbPath   = flag.String("d", "../../internal/db/test.sqlite", "Path to sqlite database")
	xmlPath  = flag.String("x", "./test.xml", "Path to nmap xml test file")
	nmapArgs = flag.String("n", "-sS localhost", "Nmap custom arguments")
)

func TestXMLParseData(t *testing.T) {
	f, err := os.Open(*xmlPath)
	if err != nil {
		t.Errorf("Error opening xml test file: %s", err)
	}

	defer f.Close()

	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		t.Errorf("Error opening database: %s", err)
	}
	defer db.Close()

	t.Run("Test XML parse normally", func(t *testing.T) {
		err = ParseData(context.Background(), f, db)
		if err != nil {
			t.Errorf("Error parsing data: %s", err)
		}
	})

	ctx, cancel := context.WithCancel(context.Background())

	t.Run("Test XML parse cancelling early", func(t *testing.T) {
		cancel()
		err = ParseData(ctx, f, db)
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

	defer stdout.Close()

	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		t.Errorf("Error opening database: %s", err)
	}
	defer db.Close()

	if err := cmd.Start(); err != nil {
		t.Errorf("Error starting nmap: %s", err)
	}

	err = ParseData(context.Background(), stdout, db)
	if err != nil {
		t.Errorf("Error parsing data: %s", err)
	}

	if err := cmd.Wait(); err != nil {
		t.Errorf("Error waiting for nmap: %s", err)
	}
}

func BenchmarkXMLParseData(b *testing.B) {
	f, err := os.Open(*xmlPath)
	if err != nil {
		b.Errorf("Error opening xml test file: %s", err)
	}

	defer f.Close()

	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		b.Errorf("Error opening database: %s", err)
	}
	defer db.Close()

	ctx := context.Background()

	for b.Loop() {
		err = ParseData(ctx, f, db)
		if err != nil {
			b.Errorf("Error parsing data: %s", err)
		}
		f.Seek(0, 0)
	}
}
