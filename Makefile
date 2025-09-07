build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o paper-tui 

dev:
	go run . --config ./test/assets/config.json

tidy:
	go mod tidy

dev-ghostty:
	ghostty --title=\"paper-tui\" -e make dev

test-all:
	go test ./... -v -cover