postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root small_bank

migrateup:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/small_bank?sslmode=disable" -verbose up

dropdb:
	docker exec -it postgres14 dropdb small_bank

migratedown:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/small_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go test -v -cover ./...

.PHONY:postgres createdb dropdb migrateup migratedown sqlc test server