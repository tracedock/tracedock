build:
	mkdir -p _build/bin
	go build \
		-ldflags \
			"-X github.com/tracedock/tracedock/cmd/tracedock/version.BuildVersion=`git describe --tags --abbrev=0 --always`" \
		-o _build/bin/tracedock \
		cmd/tracedock/main.go

test:
	mockery
	go test -v -coverprofile cover.out ./...

coverage: test
	go tool cover -html=cover.out

make develop:
	air --build.cmd "make build" --build.bin _build/bin/tracedock server start --config configs/tracedock.yaml
