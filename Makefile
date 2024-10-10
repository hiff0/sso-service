run-local:
	go run ./cmd/sso/main.go --config=./config/local.yml

migrate:
	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations