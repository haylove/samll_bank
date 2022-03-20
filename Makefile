postgres:
	docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it postgres14 createdb --username=root --owner=root small_bank

dropdb:
	docker exec -it postgres14 dropdb small_bank

migrateupa:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/small_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/small_bank?sslmode=disable" -verbose up 1

migratedowna:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/small_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/small_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

server:
	go run main.go

test:
	go test -v -cover ./...

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/haylove/small_bank/db/sqlc Store

.PHONY:postgres createdb dropdb migrateup migratedown sqlc test server mock