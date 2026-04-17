.PHONY: build run test python-sidecar up down

build:
	cd go && go build -o bin/atlas.exe .

run:
	cd go && go run .

test:
	cd go && go test ./...

python-sidecar:
	cd python && python main.py

up:
	docker compose up -d

down:
	docker compose down
