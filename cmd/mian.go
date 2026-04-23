package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/l0ng7h0r/golang/internal/handler"
	"github.com/l0ng7h0r/golang/internal/middleware"
	"github.com/l0ng7h0r/golang/internal/repository"
	"github.com/l0ng7h0r/golang/internal/usecase"
	"github.com/l0ng7h0r/golang/pkg/config"
	"github.com/l0ng7h0r/golang/pkg/database"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.DBDsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// --- Repositories ---
	userRepo     := repository.NewUserRepository(db)
	sellerRepo   := repository.NewSellerRepository(db)
	productRepo  := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	cartRepo     := repository.NewCartRepository(db)
	orderRepo    := repository.NewOrderRepository(db)
	paymentRepo  := repository.NewPaymentRepository(db)
	shipmentRepo := repository.NewShipmentRepository(db)

	// --- Usecases ---
	authUsecase     := usecase.NewAuthUsecase(userRepo)
	sellerUsecase   := usecase.NewSellerUsecase(sellerRepo, userRepo)
	productUsecase  := usecase.NewProductUsecase(productRepo, categoryRepo)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	cartUsecase     := usecase.NewCartUsecase(cartRepo, productRepo)
	orderUsecase    := usecase.NewOrderUsecase(orderRepo, cartRepo, productRepo)
	paymentUsecase  := usecase.NewPaymentUsecase(paymentRepo, orderRepo)
	shipmentUsecase := usecase.NewShipmentUsecase(shipmentRepo)

	// --- Handlers ---
	authHandler     := handler.NewAuthHandler(authUsecase)
	sellerHandler   := handler.NewSellerHandler(sellerUsecase)
	productHandler  := handler.NewProductHandler(productUsecase, sellerUsecase)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)
	cartHandler     := handler.NewCartHandler(cartUsecase)
	orderHandler    := handler.NewOrderHandler(orderUsecase)
	paymentHandler  := handler.NewPaymentHandler(paymentUsecase)
	shipmentHandler := handler.NewShipmentHandler(shipmentUsecase)

	// --- Middleware ---
	authMiddleware := middleware.NewAuthMiddleware(authUsecase)

	app := fiber.New()
	api := app.Group("/api")

	// ── Public routes ─────────────────────────────────────────────────────────
	api.Post("/register", authHandler.Register)
	api.Post("/login", authHandler.Login)
	api.Post("/refresh", authHandler.Refresh)

	// Public product/category browsing
	api.Get("/products", productHandler.GetAllProducts)
	api.Get("/products/category/:id", productHandler.GetProductsByCategory)
	api.Get("/products/seller/:id", productHandler.GetProductsBySeller)
	api.Get("/products/:id", productHandler.GetProductByID)
	api.Get("/categories", categoryHandler.GetAllCategories)
	api.Get("/categories/:id", categoryHandler.GetCategoryByID)

	// ── Authenticated user routes ──────────────────────────────────────────────
	user := api.Group("/user")
	user.Use(authMiddleware.Auth)

	// Cart
	user.Get("/cart", cartHandler.GetCart)
	user.Post("/cart/items", cartHandler.AddItem)
	user.Put("/cart/items", cartHandler.UpdateItem)
	user.Delete("/cart/items/:productId", cartHandler.RemoveItem)
	user.Delete("/cart", cartHandler.ClearCart)

	// Orders
	user.Post("/orders", orderHandler.CreateOrder)
	user.Get("/orders", orderHandler.GetMyOrders)
	user.Get("/orders/:id", orderHandler.GetOrderByID)

	// Payments
	user.Post("/payments", paymentHandler.CreatePayment)
	user.Get("/payments/order/:orderId", paymentHandler.GetPaymentByOrder)

	// Shipment tracking
	user.Get("/shipments/order/:orderId", shipmentHandler.GetShipmentByOrder)

	// ── Seller routes ──────────────────────────────────────────────────────────
	seller := api.Group("/seller")
	seller.Use(authMiddleware.Auth)
	seller.Use(authMiddleware.RequireRole("seller"))

	seller.Post("/products/create", productHandler.CreateProduct)
	seller.Put("/products/update/:id", productHandler.UpdateProduct)
	seller.Delete("/products/delete/:id", productHandler.DeleteProduct)

	// ── Admin routes ───────────────────────────────────────────────────────────
	admin := api.Group("/admin")
	admin.Use(authMiddleware.Auth)
	admin.Use(authMiddleware.RequireRole("admin"))

	// User management
	admin.Post("/users", authHandler.CreateUser)
	admin.Get("/users", authHandler.GetAllUsers)
	admin.Get("/users/:id", authHandler.GetUserByID)
	admin.Delete("/users/:id", authHandler.DeleteUser)

	// Seller management
	admin.Post("/sellers", sellerHandler.CreateSeller)
	admin.Get("/sellers", sellerHandler.GetAllSellers)
	admin.Get("/sellers/:id", sellerHandler.GetSellerByID)
	admin.Delete("/sellers/:id", sellerHandler.DeleteSeller)
	admin.Put("/sellers/:id", sellerHandler.UpdateSeller)

	// Category management
	admin.Post("/categories/create", categoryHandler.CreateCategory)
	admin.Put("/categories/update/:id", categoryHandler.UpdateCategory)
	admin.Delete("/categories/delete/:id", categoryHandler.DeleteCategory)

	// Order management
	admin.Get("/orders", orderHandler.GetAllOrders)
	admin.Patch("/orders/:id/status", orderHandler.UpdateOrderStatus)

	// Payment management
	admin.Patch("/payments/:id/confirm", paymentHandler.ConfirmPayment)

	// Shipment management
	admin.Post("/shipments/create", shipmentHandler.CreateShipment)
	admin.Patch("/shipments/:id/status", shipmentHandler.UpdateStatus)
	admin.Patch("/shipments/:id/tracking", shipmentHandler.UpdateTracking)

	app.Listen(":" + cfg.AppPort)
}