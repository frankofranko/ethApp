startdb:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_DB=eth_db -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine
	sleep 2

stopdb:
	docker stop postgres12; docker rm postgres12

connect_dbshell:
	docker exec -it postgres12 psql -U root eth_db

migrate:
	migrate -database "postgres://root:secret@localhost:5432/eth_db?sslmode=disable" -path db/migrations up
	#migrate -database "postgres://eth:password@localhost:5402/eth_db" -path db/migrations up

setup: stopdb startdb migrate

start_server:
	go run server.go

start_indexer:
	go run start_indexer.go $(block_num) $(worker_num)
