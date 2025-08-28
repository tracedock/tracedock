build:
	mkdir -p _build/bin
	go build -o _build/bin/tracedock cmd/tracedock/main.go

test:
	mockery
	go test -v -coverprofile cover.out ./...

coverage: test
	go tool cover -html=cover.out
