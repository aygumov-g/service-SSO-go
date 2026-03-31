up:
	docker compose up -d --build --force-recreate

down:
	docker compose down

logs:
	docker compose logs -f

restart:
	docker compose restart

run:
	go run cmd/app/main.go
