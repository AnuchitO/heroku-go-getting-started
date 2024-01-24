db-up:
	docker-compose up

db-down:
	docker-compose down

seed:
	ENV=LOCAL go run seeds/seed.go

run:
	ENV=LOCAL go run main.go

run-dev:
	ENV=DEV go run main.go

health:
	curl http://localhost:8080/health

# Run test
test:
	go test -v ./...

# Run test coverage
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

test-integration:
	ENV=TEST go test -v -tags integration ./...

# install swag for generate swagger file
swagger-install:
	go install github.com/swaggo/swag/cmd/swag@latest

# generate swagger file and format comment of godoc
swagger:
	swag init
	swag fmt
