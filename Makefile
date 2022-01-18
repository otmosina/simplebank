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

migrateup1:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1


migratedown:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1


sqlc:
	sqlc generate

goinit:
	go mod init github.com/otmosina/simplebank

gotidy:
	go mod tidy

test:
	richgo test -cover ./...

test-simple:
	go test -v -cover  ./...
# gotest:
# 	/usr/local/go/bin/go test  -timeout 30s -run ^main_test.go$

psql:
	docker exec -it postgres12 psql -U root -d simple_bank

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go --build_flags=--mod=mod github.com/otmosina/simplebank/db/sqlc Store

createmigration:
	migrate create -ext sql -dir db/migrations -seq add_users


.PHONY: postgres createdb dropdb migrateup migratedown