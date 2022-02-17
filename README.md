# currency-master

To run documentation locally: `~/go/bin/godoc -http=:6060` in root dir

To execute tests with coverage:
`go test ./...  -coverpkg=./... -coverprofile ./coverage.out
go tool cover -func ./coverage.out`