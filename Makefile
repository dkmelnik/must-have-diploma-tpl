up:
	go run cmd/gophermart/main.go

up.accrual:
	go run cmd/accrual/accrual_darwin_arm64

stop:
	kill -9 $(lsof -t -i :8080)

migrate.create:
	docker-compose run --rm migrate bash -c "migrate create -ext sql -dir /app/internal/db/pg/migrate ${name}"

mock:
	mockgen -source=service.go -destination=mocks/mock.go -package=mocks
