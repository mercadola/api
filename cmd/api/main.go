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
	"github.com/mercadola/api/internal/infrastruture/config"
	"github.com/mercadola/api/internal/infrastruture/database"
	"github.com/mercadola/api/internal/product"
	shoppinglist "github.com/mercadola/api/internal/shopping_list"
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

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health-check"))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(middleware.WithValue("jwt", cfg.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", cfg.JWTExpiresIn))

	customerRepository := customer.NewCustomerRepository(mongoClient, cfg, logger)
	customerEntity := &customer.Customer{}
	customerService := customer.NewService(customerRepository, logger, customerEntity)
	customerHandler := customer.NewHandler(customerService)
	customerHandler.RegisterRoutes(r)

	shoppingListRepository := shoppinglist.NewRepository(mongoClient, logger, cfg.DB, cfg.ShoppingListCollection)
	shoppingList := &shoppinglist.ShoppingList{}
	shoppinglistService := shoppinglist.NewService(shoppingListRepository, shoppingList)
	shoppinglistHandler := shoppinglist.NewHandler(shoppinglistService)
	shoppinglistHandler.RegisterRoutes(r, cfg.TokenAuth)

	productRepository := product.NewRepository(mongoClient, cfg, logger)
	productService := product.NewService(productRepository, logger)
	productHandler := product.NewHandler(productService)
	productHandler.RegisterRoutes(r, cfg.TokenAuth)

	fmt.Println("Starting server at port: " + cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatalf("Unable start a server: %v", err)
	}
}
