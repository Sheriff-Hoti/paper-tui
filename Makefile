build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hyprgo 

dev:
	go run . --config ./test/assets/config.json

tidy:
	go mod tidy