migrate:
	go run cmd/migrate/main.go

migrate-new:
	echo "please run: migrate create -ext sql -dir db/migrations create_table_xxx"

#################################### 

docker-compose:
	docker compose down -v && docker compose up

docker-validate:
	docker ps --format "{{.Names}}\t{{.Status}}"
