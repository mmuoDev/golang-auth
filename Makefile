OUTPUT = main # will be archived
VERSION = 0.1
SERVICE_NAME = golang-auth

build-local:
	go build -o $(OUTPUT) ./cmd/$(SERVICE_NAME)/main.go

run: build-local
	@echo ">> Running application ..."
	RRS_PORT=7565 \
	MONGO_DB_NAME=test \
	MONGO_URL=mongodb://localhost:27017 \
	./$(OUTPUT)