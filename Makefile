.PHONY: up down tests logs restart run

up:
	docker compose -f docker-compose.yml up -d --build --force-recreate

down:
	docker compose -f docker-compose.yml down

tests:
	docker compose -f docker-compose.yml build service-sso-app-go
	docker compose -f docker-compose.tests.yml up

logs:
	docker compose logs -f

restart:
	docker compose restart

run:
	go run cmd/app/main.go
