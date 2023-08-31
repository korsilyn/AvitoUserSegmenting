include .env
export

run:
	go run -tags migrate cmd/app/main.go

compose-up:
	docker-compose up --build -d && docker-compose logs -f

compose-down:
	docker-compose down --remove-orphans

remove-volume:
	docker volume rm pg-data

linter-golangci:
	golangci-lint run

linter-hadolint:
	git ls-files --exclude='Dockerfile*' -c --ignored | xargs hadolint

migrate-create:
	migrate create -ext sql -dir migrations 'avito_user_segmenting'

migrate-up:
#	migrate -path migrations -database '$(PG_URL_LOCALHOST)?sslmode=disable' drop -f
	migrate -path migrations -database '$(PG_URL_LOCALHOST)?sslmode=disable' up

migrate-down:
	echo "y" | migrate -path migrations -database '$(PG_URL_LOCALHOST)?sslmode=disable' down

test:
	go test -v

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

mockgen:
	mockgen -source=core/service/service.go -destination=core/mocks/servicemocks/service.go -package=servicemocks
	mockgen -source=core/repo/repo.go -destination=core/mocks/repomocks/repo.go -package=repomocks

swag:
	swag init -g internal/app/app.go --parseInternal --parseDependency
