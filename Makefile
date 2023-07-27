postgres:
	  docker rm postgres12; docker run -d --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -t postgres12 dropdb simple_bank

migrateup: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migration:
	migrate create -ext sql -dir db/migration -seq init_schema

sqlc: 
	sqlc generate

test: 
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination mock/store.go github.com/AlhajiTaibu/simplebank/sqlc Store 

.PHONY: postgres createdb dropdb migrateup migratedown migration sqlc test server mock