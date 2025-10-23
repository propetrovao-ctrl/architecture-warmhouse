package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"smarthome/db"
	"smarthome/handlers"
	"smarthome/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set up database connection
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/smarthome")
	database, err := db.New(dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer database.Close()

	log.Println("Connected to database successfully")

	// Run database migrations
	log.Println("Running database migrations...")
	if err := runMigrations(database.Pool); err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}
	log.Println("Database migrations completed successfully")

	// Initialize temperature service
	temperatureAPIURL := getEnv("TEMPERATURE_API_URL", "http://temperature-api:8081")
	temperatureService := services.NewTemperatureService(temperatureAPIURL)
	log.Printf("Temperature service initialized with API URL: %s\n", temperatureAPIURL)

	// Initialize router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API routes
	apiRoutes := router.Group("/api/v1")

	// Register sensor routes
	sensorHandler := handlers.NewSensorHandler(database, temperatureService)
	sensorHandler.RegisterRoutes(apiRoutes)

	// Start server
	srv := &http.Server{
		Addr:    getEnv("PORT", ":8080"),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited properly")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// runMigrations executes database migrations
func runMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()
	
	// Create migrations table if it doesn't exist
	createMigrationsTable := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			version VARCHAR(50) NOT NULL UNIQUE,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)
	`
	_, err := pool.Exec(ctx, createMigrationsTable)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Check if sensors table exists
	var tableExists bool
	checkTableQuery := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'sensors'
		)
	`
	err = pool.QueryRow(ctx, checkTableQuery).Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("failed to check if sensors table exists: %w", err)
	}

	if !tableExists {
		log.Println("Creating sensors table...")
		
		// Create sensors table
		createSensorsTable := `
			CREATE TABLE sensors (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				type VARCHAR(50) NOT NULL,
				location VARCHAR(100) NOT NULL,
				value FLOAT DEFAULT 0,
				unit VARCHAR(20),
				status VARCHAR(20) NOT NULL DEFAULT 'inactive',
				last_updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
				created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
			)
		`
		_, err = pool.Exec(ctx, createSensorsTable)
		if err != nil {
			return fmt.Errorf("failed to create sensors table: %w", err)
		}

		// Create indexes
		indexes := []string{
			"CREATE INDEX IF NOT EXISTS idx_sensors_type ON sensors(type)",
			"CREATE INDEX IF NOT EXISTS idx_sensors_location ON sensors(location)",
			"CREATE INDEX IF NOT EXISTS idx_sensors_status ON sensors(status)",
			"CREATE INDEX IF NOT EXISTS idx_sensors_last_updated ON sensors(last_updated)",
		}

		for _, indexSQL := range indexes {
			_, err = pool.Exec(ctx, indexSQL)
			if err != nil {
				return fmt.Errorf("failed to create index: %w", err)
			}
		}

		// Add constraints
		constraints := []string{
			"ALTER TABLE sensors ADD CONSTRAINT chk_sensors_status CHECK (status IN ('active', 'inactive', 'maintenance', 'error'))",
			"ALTER TABLE sensors ADD CONSTRAINT chk_sensors_type CHECK (type IN ('temperature', 'humidity', 'pressure', 'motion', 'light'))",
			"ALTER TABLE sensors ADD CONSTRAINT chk_sensors_value_range CHECK (value >= -50 AND value <= 100)",
			"ALTER TABLE sensors ADD CONSTRAINT uq_sensors_name_location UNIQUE (name, location)",
		}

		for _, constraintSQL := range constraints {
			_, err = pool.Exec(ctx, constraintSQL)
			if err != nil {
				log.Printf("Warning: failed to add constraint (may already exist): %v", err)
			}
		}

		// Record migration
		_, err = pool.Exec(ctx, `
			INSERT INTO migrations (version, name) 
			VALUES ('001', 'create_sensors_table')
		`)
		if err != nil {
			log.Printf("Warning: failed to record migration: %v", err)
		}

		log.Println("Sensors table created successfully")
	} else {
		log.Println("Sensors table already exists, skipping migration")
	}

	return nil
}
