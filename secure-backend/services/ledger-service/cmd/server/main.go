package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"ledger-service/internal/ledger"  // ← Import CORRECTO
	"github.com/gorilla/mux"
)

var ledgerService *ledger.Service

func init() {
	ledgerService = ledger.NewService(
		os.Getenv("JWT_SECRET"),
		os.Getenv("HMAC_SECRET"),
	)
}

func secureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		sig := r.Header.Get("X-HMAC")
		if !ledgerService.ValidateHMAC(body, sig) {
			http.Error(w, "Invalid HMAC signature", http.StatusUnauthorized)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if !ledgerService.ValidateJWT(authHeader) {
			http.Error(w, "Invalid JWT", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ledgerHandler(w http.ResponseWriter, r *http.Request) {
	var tx ledger.Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, "Invalid transaction", http.StatusBadRequest)
		return
	}

	response := ledgerService.ProcessTransaction(tx)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := ledger.HealthResponse{
		Status:       "healthy",
		Service:      "ledger",
		Version:      "1.0.0",
		TotalTxCount: 0,
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler).Methods("GET")
	
	secure := r.PathPrefix("/api").Subrouter()
	secure.Use(secureMiddleware)
	secure.HandleFunc("/transaction", ledgerHandler).Methods("POST")

	// Configuración mTLS
	caCert, _ := x509.SystemCertPool()
	customCA, _ := os.ReadFile("/certs/ca.crt")
	if customCA != nil {
		caCert.AppendCertsFromPEM(customCA)
	}

	cert, err := tls.LoadX509KeyPair("/certs/ledger.crt", "/certs/ledger.key")
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
		Addr:      ":8443",
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	log.Println("🚀 Ledger Service running on :8443")
	log.Fatal(server.ListenAndServeTLS("", ""))
}