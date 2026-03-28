.PHONY: run run-debug build clean

run:
	go run ./cmd/linux/cli/main.go

run-debug:
	go run ./cmd/linux/cli/main.go -debug

run-http:
	go run ./cmd/linux/http/main.go

run-gui:
	go run ./cmd/linux/gui/main.go

build:
	go build -o build/weather-app-cli ./cmd/linux/cli/main.go
	go build -o build/weather-app-http ./cmd/linux/http/main.go
	go build -o build/weather-app-gui ./cmd/linux/gui/main.go

clean:
	rm -rf build/*
