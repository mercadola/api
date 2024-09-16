package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/mercadola/api/internal/customer"
	"github.com/mercadola/api/internal/database"
	"github.com/mercadola/api/internal/infrastruture/config"
	"github.com/mercadola/api/internal/product"
)

func init() {
	var rootPath string
	flag.StringVar(&rootPath, "rootPath", "", "Provide project path as an absolute path")
	flag.Parse()

	if rootPath == "" {
		rootPath, _ = os.Getwd()
	}
	os.Setenv("ROOT_PATH", rootPath)

	fmt.Printf("provided path was %s\n", rootPath)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger := slog.Default()

	cfg, err := config.GetConfig()
	if err != nil {
		logger.Error("Error trying load config", err)
		os.Exit(1)
	}

	uri := cfg.Database.URI
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable.")
	}

	logger.Info("Conectando ao banco...")

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

	productRepository := product.NewProductRepository(mongoClient, cfg, logger)
	productService := product.NewService(productRepository)
	productHandler := product.NewHandler(productService)
	productHandler.RegisterRoutes(router)

	customerRepository := customer.NewCustomerRepository(mongoClient, cfg, logger)
	customerService := customer.NewService(customerRepository)
	customerHandler := customer.NewHandler(customerService)
	customerHandler.RegisterRoutes(router)

	fmt.Println("Starting server at port: " + cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, router)
	if err != nil {
		log.Fatalf("Unable start a server: %v", err)
	}
}
