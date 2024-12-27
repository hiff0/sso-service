package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	//_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable, storageHost string

	flag.StringVar(&storageHost, "storage-host", "localhost", "database host name")
	flag.StringVar(&storagePath, "storage-path", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	//if storagePath == "" {
	//	storagePath = "localhost"
	//}
	//
	//if migrationsPath == "" {
	//	panic("migrations path is required")
	//}

	//m, err := migrate.New(
	//	"file://"+migrationsPath,
	//	fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	//)
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("postgres://sso:sso@%s:5432/sso?sslmode=disable", storageHost))
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("no migrates to apply")
			return
		}

		panic(err)
	}

	log.Println("migration apply successfully")
}
