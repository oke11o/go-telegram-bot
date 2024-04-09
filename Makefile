env_example:
	cp .env.example .env

# TODO: add -X main.Version from tag
build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o tg_tournament_bot ./cmd/bot