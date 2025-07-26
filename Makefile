test:
	go test -v ./...

test-coverage:
	go test -coverpkg ./... -cover -coverprofile=coverage.out -v ./...
	go tool cover -html=coverage.out -o=coverage.html
	go tool cover -func=coverage.out -o=coverage.out
	grep -v '100.0%$$' coverage.out
