setup:
	go get -u github.com/swaggo/swag/cmd/swag
	go install github.com/swaggo/swag/cmd/swag@latest
	go get -u github.com/golang-migrate/migrate/v4/cmd/migrate
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	swag init -g ./cmd/server/main.go -o ./docs
	go get -u github.com/swaggo/gin-swagger
	go get -u github.com/swaggo/files
migrate-up:
	migrate -path pkg/database/migrations -database "postgresql://docker:password@localhost:5435/go_app_dev?sslmode=disable" up

migrate-down:
	migrate -path pkg/database/migrations -database "postgresql://docker:password@localhost:5435/go_app_dev?sslmode=disable" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir pkg/database/migrations -seq $$name



build-docker:
	docker compose build --no-cache

run-local:
	docker start dockerPostgres
	docker start dockerRedis
	docker start dockerMongo
	export REDIS_HOST=localhost
	export POSTGRES_DB=go_app_dev
	export POSTGRES_USER=docker
	export POSTGRES_PASSWORD=password
	export POSTGRES_PORT=5435
	export JWT_SECRET_KEY=ObL89O3nOSSEj6tbdHako0cXtPErzBUfq8l8o/3KD9g=INSECURE
	export API_SECRET_KEY=cJGZ8L1sDcPezjOy1zacPJZxzZxrPObm2Ggs1U0V+fE=INSECURE
	export POSTGRES_HOST=localhost
	go run cmd/server/main.go

up:
	docker compose up

down:
	docker compose down

restart:
	docker compose restart

build:
	go build -v ./...

test:
	go test -v ./... -race -cover

clean:
	docker stop ezzygo
	docker stop dockerPostgres
	docker rm ezzygo
	docker rm dockerPostgres
	docker rm dockerRedis
	docker image rm ezzygo-backend
	rm -rf .dbdata
