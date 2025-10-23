package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"smarthome/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var (
		databaseURL = flag.String("database-url", "", "Database connection URL")
		command     = flag.String("command", "up", "Migration command (up, status)")
	)
	flag.Parse()

	// Get database URL from environment or flag
	dbURL := *databaseURL
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/smarthome"
	}

	// Create database connection
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	// Create migrator
	migrator := migrations.NewMigrator(pool)

	// Execute command
	switch *command {
	case "up":
		fmt.Println("Running database migrations...")
		if err := migrator.RunMigrations(ctx); err != nil {
			log.Fatalf("Migration failed: %v\n", err)
		}
		fmt.Println("Migrations completed successfully!")

	case "status":
		fmt.Println("Migration status:")
		status, err := migrator.GetMigrationStatus(ctx)
		if err != nil {
			log.Fatalf("Failed to get migration status: %v\n", err)
		}

		if len(status) == 0 {
			fmt.Println("No migrations applied")
		} else {
			for _, s := range status {
				fmt.Printf("  %s: %s (applied at %s)\n",
					s["version"], s["name"], s["applied_at"])
			}
		}

	default:
		log.Fatalf("Unknown command: %s. Use 'up' or 'status'\n", *command)
	}
}
