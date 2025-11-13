package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"auth-service/internal/auth"  // ← Import CORRECTO
	"github.com/gorilla/mux"
)

var authService *auth.Service

func init() {
	authService = auth.NewService(
		os.Getenv("JWT_SECRET"),
		os.Getenv("HMAC_SECRET"),
	)
}

func issueTokenHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	sig := r.Header.Get("X-HMAC")
	if !authService.ValidateHMAC(body, sig) {
		http.Error(w, "Firma HMAC inválida", http.StatusUnauthorized)
		return
	}

	token, err := authService.GenerateToken("ledger")
	if err != nil {
		http.Error(w, "Error generando token", http.StatusInternalServerError)
		return
	}

	response := auth.TokenResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := auth.HealthResponse{
		Status:  "healthy",
		Service: "auth",
		Version: "1.0.0",
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/token", issueTokenHandler).Methods("POST")

	// Configuración mTLS (igual que antes)
	caCert, _ := x509.SystemCertPool()
	customCA, _ := os.ReadFile("/certs/ca.crt")
	if customCA != nil {
		caCert.AppendCertsFromPEM(customCA)
	}

	cert, err := tls.LoadX509KeyPair("/certs/auth.crt", "/certs/auth.key")
	if err != nil {
		log.Fatal("Error cargando certificados:", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCert,
		MinVersion:   tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":8442",
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	log.Println("🚀 Auth Service running on :8442")
	log.Fatal(server.ListenAndServeTLS("", ""))
}