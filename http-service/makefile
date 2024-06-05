postgres:
	docker run --rm -d --name postgres15 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=1079 -p 5432:5432 postgres:15
migrate:
	migrate create -ext sql -dir database/migration -seq init_schema
createdb:
	docker exec -it postgres15 createdb --username=postgres --owner=postgres testm

migrateup:
	migrate -path database/migration -database "postgresql://postgres:1079@localhost:5432/testm?sslmode=disable" -verbose up

migratedown:
	migrate -path  database/migration -database "postgresql://postgres:1079@localhost:5432/testm?sslmode=disable" -verbose down

dropdb:
	docker exec -it postgres15 dropdb -U postgres testm

redis:
	docker run -d -p 6379:6379 --name my-redis-container redis

protoc-email:
	protoc --go_out=. --go-grpc_out=. ./pb/email.proto

protoc-auth:
	 protoc --go_out=. --go-grpc_out=. ./pb/auth.proto
.PHONY: postgres createdb dropdb