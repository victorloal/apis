package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goHEncryption/internal/he"
	"goHEncryption/pkg/config"

	"github.com/gorilla/mux"
)

func main() {
	// Cargar configuración
	cfg := config.Load()
	
	log.Printf("Iniciando %s v%s en %s", 
		cfg.App.Name, 
		cfg.App.Version, 
		cfg.App.Environment,
	)

	// Crear router
	router := mux.NewRouter()
	
	// Middleware de logging básico
	router.Use(loggingMiddleware)
	
	// Configurar rutas
	he.SetupRoutes(router)
	
	// Ruta de health check
	router.HandleFunc("/health", healthHandler).Methods("GET")
	
	// Configurar servidor
	server := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}
	
	// Iniciar servidor en goroutine
	go func() {
		log.Printf("Servidor escuchando en %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error iniciando servidor: %v", err)
		}
	}()
	
	// Esperar señal de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Apagando servidor...")
	
	// Shutdown graceful
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error en shutdown: %v", err)
	}
	
	log.Println("Servidor apagado correctamente")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.Get()
	
	response := map[string]interface{}{
		"status":      "ok",
		"service":     cfg.App.Name,
		"version":     cfg.App.Version,
		"environment": cfg.App.Environment,
		"timestamp":   time.Now().Format(time.RFC3339),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Llamar al siguiente handler
		next.ServeHTTP(w, r)
		
		// Log después de completar la request
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}