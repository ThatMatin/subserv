package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thatmatin/subserv/internal/controller"
	"github.com/thatmatin/subserv/internal/db"
	"github.com/thatmatin/subserv/internal/middleware"
	"github.com/thatmatin/subserv/internal/repo"
	"github.com/thatmatin/subserv/internal/routers"
	"github.com/thatmatin/subserv/internal/service"
)

func RunAppandServe() {
	r := gin.Default()
	database, err := db.Setup()
	if err != nil {
		log.Fatalf("failed to setup database: %v", err)
	}

	paymentProcessor := service.NewDummyPaymentProcessor()

	productRepo := repo.NewProductRepository(database)
	userRepo := repo.NewUserRepository(database)
	subscriptionRepo := repo.NewSubscriptionRepository(database)

	productService := service.NewProductService(productRepo)
	userService := service.NewUserService(userRepo)
	subscriptionService := service.NewSubscriptionService(subscriptionRepo, productService, userService, paymentProcessor)

	productController := controller.NewProductController(&productService)
	subscriptionController := controller.NewSubscriptionController(&subscriptionService)
	routers.RegisterProductRoutes(r, productController)
	routers.RegisterSubscriptionRoutes(r, subscriptionController)

	r.Use(middleware.AuthMiddleware())

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
