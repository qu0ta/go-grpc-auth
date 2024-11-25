package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string
	flag.StringVar(&storagePath, "storage-path", "", "path for storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path for migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrating table")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is empty")
	}
	if migrationsPath == "" {
		panic("migrations-path is empty")
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("Failed to migrate: ", err)
			return
		} else {
			panic(err)
		}
	}
	fmt.Println("Successfully migrated")
}
