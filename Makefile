make-pg:
	docker run --name pg-container -p 5433:5432 -e POSTGRES_USER=$(E) -e POSTGRES_PASSWORD=$(E) -d postgres:13.6

start-pg:
	docker start pg-container

createdb:
	docker exec -it pg-container createdb --username=$(E) --owner=$(E) bank

dropdb:
	docker exec -it pg-container dropdb --username=$(E) bank

migrateup:
	migrate -path db/migration -database "postgresql://$(E):$(E)@localhost:5433/bank?sslmode=disable" --verbose up $(V)

migratedown:
	migrate -path db/migration -database "postgresql://$(E):$(E)@localhost:5433/bank?sslmode=disable" --verbose down $(V)

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: make-pg start-pg createdb dropdb migrateup migratedown sqlc test server