createdb:
	docker exec -it postgres12 createdb --user=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

postgres:
	docker run -p 5432:5432 --name postgres12 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

postgres-start:
	docker start postgres12


migrateup:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

goinit:
	go mod init github.com/otmosina/simplebank

gotidy:
	go mod tidy

test:
	go test -v -cover  ./...
# gotest:
# 	/usr/local/go/bin/go test  -timeout 30s -run ^main_test.go$

psql:
	docker exec -it postgres12 psql -U root -d simple_bank

.PHONY: postgres createdb dropdb migrateup migratedown