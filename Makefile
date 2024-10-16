run-local:
	go run ./cmd/sso/main.go --config=./config/local.yml

test:
	go test ./tests -v

test-cover:
	go test ./tests -v -cover

migrate:
	go run ./cmd/migrator/main.go --storage-path=./storage/sso.db --migrations-path=./migrations

test-migrate:
	go run ./cmd/migrator/main.go --storage-path=./storage/sso.db --migrations-path=./tests/migrations --migrations-table=migrations_test