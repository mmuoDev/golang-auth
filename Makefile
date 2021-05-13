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
	JWT_ACCESS_SECRET=T52N6pRxNZDW45UR \
	JWT_REFRESH_SECRET=Q768EuNprKx4uhGj \
	REDIS_DSN=localhost:6379 \
	./$(OUTPUT)