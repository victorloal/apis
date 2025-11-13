package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/google/uuid"

	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/multiparty"
	"github.com/tuneinsight/lattigo/v6/schemes/bgv"
	"github.com/tuneinsight/lattigo/v6/utils/sampling"
)

// --- Copiar tus tipos y funciones Lattigo necesarias (simplificadas) ---
type Party struct {
	id       int
	sk       *rlwe.SecretKey
	shamirPt multiparty.ShamirPublicPoint
	tkShare  multiparty.ShamirSecretShare
	combiner multiparty.Combiner
}

var (
	// Estado global (protected por mutex)
	params     bgv.Parameters
	parties    []*Party
	publicKey  *rlwe.PublicKey
	stateMutex sync.RWMutex

	// Store ciphertexts in memory with UUIDs
	ctStore   = map[string]*rlwe.Ciphertext{}
	ctMutex   sync.RWMutex
	apiKeyEnv = "MY_API_KEY" // env var name expected
)

// ------------------ UTIL / LATTIGO SETUP (adaptado de tu código) ------------------

func initParams() bgv.Parameters {
	params, err := bgv.NewParametersFromLiteral(bgv.ParametersLiteral{
		LogN:             15,
		LogQ:             []int{54, 54, 54, 54},
		LogP:             []int{55},
		PlaintextModulus: 1999699969,
	})
	if err != nil {
		log.Fatal(err)
	}
	return params
}

func createParties(N, t int, params bgv.Parameters) []*Party {
	gen := rlwe.NewKeyGenerator(params)
	parties := make([]*Party, N)
	shamirPts := make([]multiparty.ShamirPublicPoint, N)

	for i := 0; i < N; i++ {
		parties[i] = &Party{
			id:       i,
			sk:       gen.GenSecretKeyNew(),
			shamirPt: multiparty.ShamirPublicPoint(i + 1),
		}
		shamirPts[i] = parties[i].shamirPt
	}

	thresholdizer := multiparty.NewThresholdizer(params)

	for i := range parties {
		poly, err := thresholdizer.GenShamirPolynomial(t, parties[i].sk)
		if err != nil {
			log.Fatal(err)
		}

		for j := range parties {
			share := thresholdizer.AllocateThresholdSecretShare()
			thresholdizer.GenShamirSecretShare(parties[j].shamirPt, poly, &share)

			if i == 0 {
				parties[j].tkShare = share
			} else {
				thresholdizer.AggregateShares(parties[j].tkShare, share, &parties[j].tkShare)
			}
		}
	}

	for i := range parties {
		parties[i].combiner = multiparty.NewCombiner(*params.GetRLWEParameters(), parties[i].shamirPt, shamirPts, t)
	}
	return parties
}

func getShamirPoints(parties []*Party) []multiparty.ShamirPublicPoint {
	pts := make([]multiparty.ShamirPublicPoint, len(parties))
	for i, p := range parties {
		pts[i] = p.shamirPt
	}
	return pts
}

func collectivePublicKey(parties []*Party, t int, params bgv.Parameters) *rlwe.PublicKey {
	crs, err := sampling.NewPRNG()
	if err != nil {
		log.Fatalf("error creando PRNG: %v", err)
	}

	ckg := multiparty.NewPublicKeyGenProtocol(params)
	crp := ckg.SampleCRP(crs)

	activeParties := parties[:t]
	shamirPts := getShamirPoints(activeParties)

	shares := make([]multiparty.PublicKeyGenShare, len(activeParties))
	tsks := make([]*rlwe.SecretKey, len(activeParties))

	for i, p := range activeParties {
		shares[i] = ckg.AllocateShare()
		tsks[i] = rlwe.NewSecretKey(params)

		err := p.combiner.GenAdditiveShare(shamirPts, p.shamirPt, p.tkShare, tsks[i])
		if err != nil {
			log.Fatal(err)
		}

		ckg.GenShare(tsks[i], crp, &shares[i])
	}

	aggShare := shares[0]
	for i := 1; i < len(activeParties); i++ {
		ckg.AggregateShares(aggShare, shares[i], &aggShare)
	}

	publicKey := rlwe.NewPublicKey(params)
	ckg.GenPublicKey(aggShare, crp, publicKey)

	return publicKey
}

func encryptMessage(message uint64, pk *rlwe.PublicKey, params bgv.Parameters) *rlwe.Ciphertext {
	encoder := bgv.NewEncoder(params)
	encryptor := rlwe.NewEncryptor(params, pk)

	modulus := uint64(params.PlaintextModulus())
	if message >= modulus {
		log.Fatalf("Número %d excede el módulo %d", message, modulus)
	}

	vector := make([]uint64, params.N())
	for i := range vector {
		vector[i] = message
	}

	pt := bgv.NewPlaintext(params, params.MaxLevel())
	encoder.Encode(vector, pt)

	ct, err := encryptor.EncryptNew(pt)
	if err != nil {
		log.Fatalf("error cifrando el mensaje: %v", err)
	}
	return ct
}

func sumarCiphertexts(params bgv.Parameters, ciphertexts []*rlwe.Ciphertext) *rlwe.Ciphertext {
	if len(ciphertexts) == 0 {
		return nil
	}
	evaluator := bgv.NewEvaluator(params, nil)
	resultado := ciphertexts[0].CopyNew()
	for i := 1; i < len(ciphertexts); i++ {
		if err := evaluator.Add(resultado, ciphertexts[i], resultado); err != nil {
			log.Fatal("Error sumando ciphertexts:", err)
		}
	}
	return resultado
}

func decryptCollective(ct *rlwe.Ciphertext, parties []*Party, t int, params bgv.Parameters) uint64 {
	if len(parties) < t {
		log.Fatalf("No hay suficientes partes (%d) para el umbral t=%d", len(parties), t)
	}

	collectiveSK := rlwe.NewSecretKey(params)
	shamirPts := getShamirPoints(parties)

	for _, p := range parties {
		skShare := rlwe.NewSecretKey(params)
		err := p.combiner.GenAdditiveShare(shamirPts, p.shamirPt, p.tkShare, skShare)
		if err != nil {
			log.Fatal(err)
		}
		params.RingQP().Add(collectiveSK.Value, skShare.Value, collectiveSK.Value)
	}

	decryptor := rlwe.NewDecryptor(params, collectiveSK)
	pt := bgv.NewPlaintext(params, params.MaxLevel())
	decryptor.Decrypt(ct, pt)

	encoder := bgv.NewEncoder(params)
	decoded := make([]uint64, params.N())
	encoder.Decode(pt, decoded)

	return decoded[0]
}

// ------------------ HTTP / GORILLA HANDLERS ------------------

type InitRequest struct {
	N int `json:"n"`
	T int `json:"t"`
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	var req InitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.N <= 0 || req.T <= 0 || req.T > req.N {
		http.Error(w, "invalid parameters", http.StatusBadRequest)
		return
	}

	// Inicializar params y parties
	stateMutex.Lock()
	defer stateMutex.Unlock()

	params = initParams()
	parties = createParties(req.N, req.T, params)
	publicKey = collectivePublicKey(parties, req.T, params)

	// limpiar store
	ctMutex.Lock()
	ctStore = map[string]*rlwe.Ciphertext{}
	ctMutex.Unlock()

	resp := map[string]interface{}{
		"message":        "initialized",
		"n":              req.N,
		"t":              req.T,
		"plaintext_mod":  params.PlaintextModulus(),
		"valid_range_to": params.PlaintextModulus() - 1,
	}
	json.NewEncoder(w).Encode(resp)
}

type EncryptRequest struct {
	Message uint64 `json:"message"`
}

type EncryptResponse struct {
	CipherID string `json:"cipher_id"`
}

func encryptHandler(w http.ResponseWriter, r *http.Request) {
	stateMutex.RLock()
	pk := publicKey
	paramsLocal := params
	stateMutex.RUnlock()

	if pk == nil {
		http.Error(w, "server not initialized", http.StatusBadRequest)
		return
	}

	var req EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	ct := encryptMessage(req.Message, pk, paramsLocal)

	id := uuid.New().String()
	ctMutex.Lock()
	ctStore[id] = ct
	ctMutex.Unlock()

	json.NewEncoder(w).Encode(EncryptResponse{CipherID: id})
}

type SumRequest struct {
	IDs []string `json:"ids"`
}

type SumResponse struct {
	SumID string `json:"sum_id"`
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	stateMutex.RLock()
	paramsLocal := params
	stateMutex.RUnlock()

	var req SumRequest
	if len(req.IDs) == 0 {
		// decode body even if empty to validate
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.IDs) == 0 {
		http.Error(w, "no ids provided", http.StatusBadRequest)
		return
	}

	// recuperar ciphertexts
	ctMutex.RLock()
	ctList := make([]*rlwe.Ciphertext, 0, len(req.IDs))
	for _, id := range req.IDs {
		ct, ok := ctStore[id]
		if !ok {
			ctMutex.RUnlock()
			http.Error(w, fmt.Sprintf("cipher not found: %s", id), http.StatusNotFound)
			return
		}
		ctList = append(ctList, ct)
	}
	ctMutex.RUnlock()

	sumCt := sumarCiphertexts(paramsLocal, ctList)
	if sumCt == nil {
		http.Error(w, "sum failed", http.StatusInternalServerError)
		return
	}

	sumID := uuid.New().String()
	ctMutex.Lock()
	ctStore[sumID] = sumCt
	ctMutex.Unlock()

	json.NewEncoder(w).Encode(SumResponse{SumID: sumID})
}

type DecryptRequest struct {
	CipherID  string `json:"cipher_id"`
	Threshold int    `json:"threshold"` // cuántas partes usar (t)
}

type DecryptResponse struct {
	Plaintext uint64 `json:"plaintext"`
}

func decryptHandler(w http.ResponseWriter, r *http.Request) {
	var req DecryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// obtener ciphertext
	ctMutex.RLock()
	ct, ok := ctStore[req.CipherID]
	ctMutex.RUnlock()
	if !ok || ct == nil {
		http.Error(w, "cipher not found", http.StatusNotFound)
		return
	}

	stateMutex.RLock()
	paramsLocal := params
	// usar las primeras req.Threshold parties (por simplicidad)
	if len(parties) < req.Threshold || req.Threshold <= 0 {
		stateMutex.RUnlock()
		http.Error(w, "insufficient parties or invalid threshold", http.StatusBadRequest)
		return
	}
	selected := parties[:req.Threshold]
	stateMutex.RUnlock()

	plain := decryptCollective(ct, selected, req.Threshold, paramsLocal)

	json.NewEncoder(w).Encode(DecryptResponse{Plaintext: plain})
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	stateMutex.RLock()
	defer stateMutex.RUnlock()
	ctMutex.RLock()
	defer ctMutex.RUnlock()

	resp := map[string]interface{}{
		"initialized": publicKey != nil,
		"n_parties":   len(parties),
		"stored_cts":  len(ctStore),
	}
	json.NewEncoder(w).Encode(resp)
}

// ------------------ AUTH MIDDLEWARE ------------------

func requireAPIKey(next http.Handler) http.Handler {
	expected := os.Getenv(apiKeyEnv)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if expected == "" {
			// For safety: if not set, deny
			http.Error(w, "server misconfigured: API key not set", http.StatusInternalServerError)
			return
		}
		if key != expected {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ------------------ MAIN / ROUTER / TLS ------------------

func main() {
	// Para desarrollo, setea env MY_API_KEY antes de ejecutar:
	// export MY_API_KEY="mi_super_secreta_api_key"
	if os.Getenv(apiKeyEnv) == "" {
		log.Println("WARNING: API key not set. Set environment variable", apiKeyEnv, "before production.")
	}

	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/init", initHandler).Methods("POST")
	r.HandleFunc("/encrypt", encryptHandler).Methods("POST")
	r.HandleFunc("/sum", sumHandler).Methods("POST")
	r.HandleFunc("/decrypt", decryptHandler).Methods("POST")
	r.HandleFunc("/status", statusHandler).Methods("GET")

	// Apply middleware to mutating endpoints and encrypt/decrypt
	secure := r.NewRoute().Subrouter()
	secure.Use(requireAPIKey)
	secure.HandleFunc("/init", initHandler).Methods("POST")
	secure.HandleFunc("/encrypt", encryptHandler).Methods("POST")
	secure.HandleFunc("/sum", sumHandler).Methods("POST")
	secure.HandleFunc("/decrypt", decryptHandler).Methods("POST")

	// Set up TLS server
	certFile := "cert.pem"
	keyFile := "key.pem"
	addr := ":8443"

	// If TLS cert not found, show a friendly message.
	if _, err := os.Stat(certFile); errors.Is(err, os.ErrNotExist) {
		log.Printf("TLS certificate not found (%s). Generate a self-signed for dev or provide cert/key.\n", certFile)
		log.Printf("You can generate a dev cert with: openssl req -x509 -newkey rsa:4096 -nodes -sha256 -days 365 -keyout %s -out %s\n", keyFile, certFile)
	}

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Println("Starting server on https://localhost" + addr)
	if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatal(err)
	}
}
