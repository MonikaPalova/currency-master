# currency-master

To run documentation locally: `godoc -http=:6060` in root dir

To execute tests with coverage:
`
 go test ./...   -coverpkg=./... -coverprofile coverage.out
 go tool cover -func coverage.out
`

Prerequisites:

- Mysql started on port 3306


Improvements:

- add config vars into environment variables