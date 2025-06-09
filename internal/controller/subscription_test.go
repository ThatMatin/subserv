package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	mocklib "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thatmatin/subserv/internal/middleware"
	"github.com/thatmatin/subserv/internal/mock"
	"github.com/thatmatin/subserv/internal/model"
	"github.com/thatmatin/subserv/internal/service"
)

var fixedTime = time.Date(2020, time.May, 0, 0, 0, 0, 0, time.UTC)

func TestSubscriptionController(t *testing.T) {
	router := gin.Default()

	mockSubscriptionRepo := new(mock.MockSubscriptionRepo)
	mockUserRepo := new(mock.MockUserRepo)
	paymentProcessor := service.NewDummyPaymentProcessor()
	mockProductRepo := new(mock.MockProductRepo)
	productService := service.NewProductService(mockProductRepo)
	mockSubscriptionService := service.NewSubscriptionService(mockSubscriptionRepo, productService, mockUserRepo, paymentProcessor)
	subscriptionController := NewSubscriptionController(&mockSubscriptionService)

	router.Use(middleware.AuthMiddleware())
	router.GET("/subscriptions/:id", subscriptionController.GetSubscriptionByID)
	router.POST("/subscriptions", subscriptionController.CreateSubscription)
	router.POST("/subscriptions/:id/purchase", subscriptionController.Purchase)
	router.PATCH("/subscriptions/:id/pause", subscriptionController.PauseSubscription)
	router.PATCH("/subscriptions/:id/unpause", subscriptionController.UnpauseSubscription)
	router.PATCH("/subscriptions/:id/cancel", subscriptionController.CancelSubscription)

	t.Run("get subscription by id", func(t *testing.T) {
		mockSubscriptionRepo.On("GetByID", mocklib.Anything, uint(1)).
			Return(&model.Subscription{ID: 1, UserID: 1, ProductID: 1, State: model.Active}, nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/subscriptions/1", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), `"id":1,"user_id":1,"product_id":1,"state":1`)
		mockSubscriptionRepo.AssertExpectations(t)

		mockSubscriptionRepo.ExpectedCalls = nil
	})

	t.Run("create subscription", func(t *testing.T) {
		mockProductRepo.On("GetByID", mocklib.Anything, uint(1)).Return(&model.Product{ID: 1, Name: "Test Product", Price: 1000, TaxRate: 20, Duration: 30}, nil)
		mockUserRepo.On("Exists", mocklib.Anything, uint(1)).Return(true, nil)
		mockSubscriptionRepo.On("Create", mocklib.Anything, mocklib.Anything).Return(nil)

		w := httptest.NewRecorder()
		jsonBody := `{"product_id": 1}`
		req := httptest.NewRequest(http.MethodPost, "/subscriptions", strings.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer test-token")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)
		require.Contains(t, w.Body.String(), `"id":0,"user_id":1,"product_id":1,"state":0`)
		mockSubscriptionRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)

		mockSubscriptionRepo.ExpectedCalls = nil
		mockProductRepo.ExpectedCalls = nil
		mockUserRepo.ExpectedCalls = nil
	})

	t.Run("purchase subscription", func(t *testing.T) {
		mockSubscriptionRepo.On("GetByID", mocklib.Anything, uint(1)).
			Return(&model.Subscription{ID: 1, UserID: 1, ProductID: 1, State: model.Pending}, nil)
		mockSubscriptionRepo.On("Save", mocklib.Anything, mocklib.Anything).Return(nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/subscriptions/1/purchase", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), `Subscription purchased successfully`)
		mockSubscriptionRepo.AssertExpectations(t)
		mockSubscriptionRepo.ExpectedCalls = nil
	})

	t.Run("pause subscription", func(t *testing.T) {
		mockSubscriptionRepo.On("GetByID", mocklib.Anything, uint(1)).
			Return(&model.Subscription{ID: 1, UserID: 1, ProductID: 1, State: model.Active}, nil)
		mockSubscriptionRepo.On("Save", mocklib.Anything, mocklib.Anything).Return(nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/subscriptions/1/pause", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusAccepted, w.Code)
		mockSubscriptionRepo.AssertExpectations(t)
		mockSubscriptionRepo.ExpectedCalls = nil
	})

	t.Run("unpause subscription", func(t *testing.T) {
		mockSubscriptionRepo.On("GetByID", mocklib.Anything, uint(1)).
			Return(&model.Subscription{ID: 1, UserID: 1, ProductID: 1, State: model.Paused, PausedAt: &fixedTime}, nil)
		mockSubscriptionRepo.On("Save", mocklib.Anything, mocklib.Anything).Return(nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/subscriptions/1/unpause", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusAccepted, w.Code)
		mockSubscriptionRepo.AssertExpectations(t)
		mockSubscriptionRepo.ExpectedCalls = nil
	})

	t.Run("cancel subscription", func(t *testing.T) {
		mockSubscriptionRepo.On("GetByID", mocklib.Anything, uint(1)).
			Return(&model.Subscription{ID: 1, UserID: 1, ProductID: 1, State: model.Active}, nil)
		mockSubscriptionRepo.On("Save", mocklib.Anything, mocklib.Anything).Return(nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/subscriptions/1/cancel", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusAccepted, w.Code)
		mockSubscriptionRepo.AssertExpectations(t)
		mockSubscriptionRepo.ExpectedCalls = nil
	})
}
