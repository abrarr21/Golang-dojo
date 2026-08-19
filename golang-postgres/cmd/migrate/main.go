package main

import (
	"fmt"
	"golang-postgres/internal/config"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: make migrate <up | down>")
	}

	cfg := config.Load()

	m, err := migrate.New(
		"file://migrations",
		cfg.DbURL,
	)
	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}

	switch os.Args[1] {
	case "up":
		log.Println("Migration UP called")
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}

	case "down":
		log.Println("Migration DOWN called")
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}

	fmt.Println("Running migrations...")

}
