package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/l0ng7h0r/golang/internal/handler"
	"github.com/l0ng7h0r/golang/internal/repository"
	"github.com/l0ng7h0r/golang/internal/usecase"
	"github.com/l0ng7h0r/golang/pkg/config"
	"github.com/l0ng7h0r/golang/pkg/database"
	"github.com/l0ng7h0r/golang/internal/middleware"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.DBDsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}   

	//User routes

	userRepo := repository.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepo)
	authHandler := handler.NewAuthHandler(authUsecase)
	authMiddleware := middleware.NewAuthMiddleware(authUsecase)

	//Seller routes

	sellerRepo := repository.NewSellerRepository(db)
	sellerUsecase := usecase.NewSellerUsecase(sellerRepo, userRepo)
	sellerHandler := handler.NewSellerHandler(sellerUsecase)

	//Product routes

	productRepo := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)
	productHandler := handler.NewProductHandler(productUsecase, sellerUsecase)

	app := fiber.New()

	api := app.Group("/api")

	//Public routes
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)
	api.Post("/refresh", authHandler.Refresh)

	//Product routes
	product := api.Group("/product")
	product.Use(authMiddleware.Auth)
	product.Use(authMiddleware.RequireRole("seller"))
	product.Post("/create", productHandler.CreateProduct)
	product.Get("/products", productHandler.GetAllProducts)
	product.Get("/products/:id", productHandler.GetProductByID)
	product.Delete("/products/:id", productHandler.DeleteProduct)
	product.Put("/products/:id", productHandler.UpdateProduct)

	//Admin routes
	admin := api.Group("/admin")
	admin.Use(authMiddleware.Auth)
	admin.Use(authMiddleware.RequireRole("admin"))
	admin.Post("/create-user", authHandler.CreateUser)
	admin.Get("/users", authHandler.GetAllUsers)
	admin.Get("/users/:id", authHandler.GetUserByID)
	admin.Delete("/users/:id", authHandler.DeleteUser)
	admin.Post("/create-seller", sellerHandler.CreateSeller)
	admin.Get("/sellers", sellerHandler.GetAllSellers)
	admin.Get("/sellers/:id", sellerHandler.GetSellerByID)
	admin.Delete("/sellers/:id", sellerHandler.DeleteSeller)
	admin.Put("/sellers/:id", sellerHandler.UpdateSeller)

	app.Listen(":" + cfg.AppPort)
}