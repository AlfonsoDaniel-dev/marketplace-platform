DB_USER = poncho
DB_PASSWORD = SecurePassword
DB_NAME = shopDB
DB_CONTAINER_NAME = ShopperiaPostgres
CONTAINER_INTERNAL_PORT = 5432
CONTAINER_EXTERNAL_PORT = 3001

DATABASE:
	docker run --name $(DB_CONTAINER_NAME) -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -e POSTGRES_DB=$(DB_NAME) -p $(CONTAINER_EXTERNAL_PORT):$(CONTAINER_INTERNAL_PORT) -d postgres:latest
STOP:
	docker stop $(DB_CONTAINER_NAME)
START:
	docker start $(DB_CONTAINER_NAME)
DROP:
	docker rm $(DB_CONTAINER_NAME)
RUN:
	go run ./cmd .
BUILD:
	go build -o shopperia ./cmd/main.go
REMOVE_BINARY:
	rm shopperia