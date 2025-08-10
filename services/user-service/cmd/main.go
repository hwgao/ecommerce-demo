package main

import (
    "log"
    "net/http"
    "os"
    "github.com/gorilla/mux"
    
    "github.com/prometheus/client_golang/prometheus/promhttp"
    userHandler "ecommerce/services/user-service/internal/handler/http"
    "ecommerce/services/user-service/internal/repository/postgres"
    "ecommerce/services/user-service/internal/service"
    "ecommerce/shared/pkg/cache"
    "ecommerce/shared/pkg/database"
    "ecommerce/shared/pkg/events"
    "ecommerce/shared/pkg/middleware"
    "ecommerce/shared/pkg/metrics"
)

func main() {
    // Database connection
    db, err := database.NewPostgresConnection(os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Redis cache
    redisCache := cache.NewRedisCache(os.Getenv("REDIS_URL"))

    // Kafka event bus
    eventBus := events.NewKafkaEventBus(os.Getenv("KAFKA_BROKERS"))

    // Initialize repositories and services
    userRepo := postgres.NewUserRepository(db)
    userService := service.NewUserService(userRepo, redisCache, eventBus, os.Getenv("JWT_SECRET"))

    // Initialize handlers
    userHandler := userHandler.NewUserHandler(userService)

    // Setup router
    router := mux.NewRouter()
    
    // Middleware
    router.Use(middleware.Logging)
    router.Use(middleware.CORS)
    router.Use(middleware.Metrics)

    // Routes
    api := router.PathPrefix("/api/v1").Subrouter()
    userHandler.SetupRoutes(api.PathPrefix("/users").Subrouter())

    // Health check
    router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // Metrics endpoint
    router.Handle("/metrics", promhttp.Handler())

    // Initialize metrics
    metrics.Init("user_service")

    log.Println("User service starting on port 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}
