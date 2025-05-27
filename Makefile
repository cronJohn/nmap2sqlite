TEST_DB = internal/db/test.sqlite

gen:
	sqlc generate

gen-alt:
	sqlc generate -f ".other.sqlc.yaml"

test-xml: build-db
	cd ./pkg/nmap && go test -v -run TestXMLParseData

bench-xml: build-db
	cd ./pkg/nmap && go test -bench=BenchmarkXMLParseData -run ==

test-nmap: build-db
	cd ./pkg/nmap && sudo go test -v -run TestNmapScanParseData

build-db:
	rm -f $(TEST_DB)
	sqlite3 $(TEST_DB) < "internal/db/schema.sql"

ro-custom-sql:
	chmod 444 pkg/db/custom.sql.go

rw-custom-sql:
	chmod 777 pkg/db/custom.sql.go
