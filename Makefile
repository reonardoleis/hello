build-server:
	go build -o bin/server cmd/server/server.go

build-client:
	go build -o bin/client cmd/client/client.go

run-server: 
	go run cmd/server/server.go

run-client:
	go run cmd/client/client.go