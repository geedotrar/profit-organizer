package main

import (
	"laba_service/config"
	"laba_service/internal/handlers"
	repository "laba_service/internal/repositories"
	"laba_service/internal/routes"
	service "laba_service/internal/services"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	server()
}

func server() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	g := gin.Default()
	g.Use(gin.Recovery())

	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	gormConfig := config.NewGormPostgres()
	if gormConfig == nil {
		log.Fatal("Failed to initialize database connection")
	}

	db := gormConfig.GetConnection()
	if db == nil {
		log.Fatal("Database connection is nil")
	}

	labaGroup := g.Group("/laba")
	labaRepo := repository.NewLabaRepository(gormConfig)
	labaSvc := service.NewLabaService(labaRepo)
	labaHdl := handlers.NewproductHandler(labaSvc)
	labaRouter := routes.NewProductRouter(labaGroup, labaHdl)
	labaRouter.Mount()
	g.Static("/uploads", "./uploads")

	g.Run(":8080")
}
