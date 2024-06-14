# go install github.com/pressly/goose/v3/cmd/goose@latest
generate_sql_migration:
	goose -dir migrations create $(name) sql

apply_migrations:
	goose -dir migrations up

vet:
	go vet ./...

build_app:
	go build -o ./bin/app ./cmd/app/main.go

build_worker:
	go build -o ./bin/worker ./cmd/worker/main.go
