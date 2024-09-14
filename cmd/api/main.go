package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/mercadola/api/internal/database"
	"github.com/mercadola/api/internal/product"
)

func main() {
	logger := slog.Default()

	if os.Getenv("ENV") != "CONTAINER" {
		err := godotenv.Load()
		if err != nil {
			logger.Error("Error loading .env file")
		}
	}
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable.")
	}
	logger.Info("Conectando ao banco...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := database.NewClient(uri, ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	logger.Info("Conexão realizada com sucesso...")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/health-check"))
	router.Use(middleware.Timeout(60 * time.Second))

	productRepository := product.NewProductRepository(mongoClient, logger)
	productService := product.NewService(productRepository)
	productHandler := product.NewHandler(productService)

	router.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.Find)
	})

	fmt.Println("Starting server at port: " + os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), router)
	if err != nil {
		log.Fatalf("Unable start a server: %v", err)
	}
}
