package main

import (
	"log"
	"os"

	"payslip-system/internal/config"
	"payslip-system/internal/controllers/api"
	"payslip-system/internal/database"
	"payslip-system/internal/middleware"
	"payslip-system/internal/providers"
	"payslip-system/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run migrations
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	services := providers.NewServices(repos)

	// Initialize Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.RequestLogger())

	// Setup routes with repositories (fixed routing)
	api.SetupRoutesWithRepos(r, services, repos)

	// Seed database if needed
	if cfg.SeedDatabase {
		if err := database.SeedDatabase(repos); err != nil {
			log.Printf("Warning: Failed to seed database: %v", err)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Default admin credentials: username=admin, password=admin123")
	log.Printf("Employee credentials: username=employee1-100, password=employee1-100")

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
