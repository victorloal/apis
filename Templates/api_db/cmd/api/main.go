package main

import (
    "api_db/internal/routes"
    "api_db/pkg/config"
    "api_db/pkg/db"
    "log"
    "net/http"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig("./configs")
    if err != nil {
        log.Fatal("Cannot load config:", err)
    }

    // Connect to database
    database, err := db.ConnectDB(&cfg.Database)
    if err != nil {
        log.Fatal("Cannot connect to database:", err)
    }

    // Setup routes
    router := routes.SetupRoutes(database)

    // Start server
    log.Printf("Server starting on port %s", cfg.Server.Port)
    log.Fatal(http.ListenAndServe(cfg.Server.Port, router))
}