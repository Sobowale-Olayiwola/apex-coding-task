sqlc:
	sqlc generate

postgres:
	docker run --name simple_dice_postgres -p 5435:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it simple_dice_postgres createdb --username=root --owner=root simple_dice

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5435/simple_dice?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5435/simple_dice?sslmode=disable" -verbose down

dropdb:
	docker exec -it postgres12tut dropdb simple_dice_postgres

server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown mock 