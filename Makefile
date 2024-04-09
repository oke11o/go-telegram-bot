env_example:
	cp .env.example .env

# TODO: add -X main.Version from tag
build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X 'main.Version=v0.1.0'" -o tg_tournament_bot ./cmd/bot

server_build:
	git describe --tags `git rev-list --tags --max-count=1`
	CGO_ENABLED=1 go build -ldflags "-s -w -X 'main.Version=v0.1.0'" -o bin/tg_tournament_bot ./cmd/bot