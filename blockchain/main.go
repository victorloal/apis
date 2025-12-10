package main

import (
	"encoding/json"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"
	"io/ioutil"
    "github.com/rs/cors"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "/home/victor/Documentos/blockcahin-votation/acme-network/crypto-config/org1.acme.com"
	certPath     = cryptoPath + "/users/admin@org1.acme.com/msp/signcerts/"
	keyPath      = cryptoPath + "/users/admin@org1.acme.com/msp/keystore/"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.acme.com/tls/ca.crt"
	peerEndpoint = "dns:///localhost:7051"
	gatewayPeer  = "peer0.org1.acme.com"
)

func main() {
    
    http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./public"))))

    http.HandleFunc("/newCandidato", newCandidateHandler) 
    http.HandleFunc("/newVoto", newVoteHandler)
    http.HandleFunc("/newElection", newElectionHandler)
    http.HandleFunc("/getCandidates", getCandidatesHandler)
    http.HandleFunc("/getVotes", getVotesHandler)
    http.HandleFunc("/getElection", getElectionHandler)

    
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"}, // Permitir todas las solicitudes de origen
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"}, // Permitir GET, POST y OPTIONS
        AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Permitir estos encabezados
        AllowCredentials: true,
    })

    // Iniciar el servidor en el puerto 8088 con CORS habilitado
    fmt.Println("Servidor iniciado en http://localhost:3002")
    handler := c.Handler(http.DefaultServeMux) // Agregar el middleware CORS

    if err := http.ListenAndServe(":3002", handler); err != nil {
        fmt.Printf("Error al iniciar el servidor: %v\n", err)
    }
}


func newElectionHandler(w http.ResponseWriter, r *http.Request) {
	
	var body struct {
    	Ui string `json:"ui"`;
		Nombre string `json:"nombre"`;
		Inicio time.Time `json:"inicio"`;
		Fin time.Time `json:"fin"`;
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	println(body.Ui)
	println(body.Nombre)
	println(body.Inicio.String())
	println(body.Fin.String())

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error de conexión con el Gateway: %v", err), http.StatusInternalServerError)
		return
	}
	defer gw.Close()

	// Configuración del contrato
	chaincodeName := "votacion"
	channelName := "marketplace"
	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	inicio := body.Inicio.Format(time.RFC3339)
	fin := body.Fin.Format(time.RFC3339)
	
	// Consultar los votos de un candidato
	result, err := newElection(contract, body.Ui, body.Nombre, inicio, fin) // Usar el nombre dinámico
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener a los votantes: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"elección": result})
}

func newElection (contract *client.Contract, ui string, nombre string, inicio string, fin string) (interface{}, error) {

	
	evaluateResult, err := contract.SubmitTransaction("NuevaVotacion", ui, nombre, inicio, fin)
	if err != nil {
		return nil, err
	}
	println(1)
	return string(evaluateResult),nil
}	

func newCandidateHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
    	UiVotacion string `json:"uiVotacion"`;
		UiCandidato string `json:"uiCandidato"`;
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	println(body.UiVotacion)
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error de conexión con el Gateway: %v", err), http.StatusInternalServerError)
		return
	}
	defer gw.Close()

	// Configuración del contrato
	chaincodeName := "votacion"
	channelName := "marketplace"
	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	// Consultar los votos de un candidato
	result, err := newCandidate(contract, body.UiVotacion, body.UiCandidato) // Usar el nombre dinámico
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener a los votantes: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"elección": result})
}

func newCandidate (contract *client.Contract, uiVotacion string, idCandidate string) (interface{}, error) {
	evaluateResult, err := contract.SubmitTransaction("AgregarCandidato", uiVotacion, idCandidate)
	if err != nil {
		return nil, err
	}
	println(uiVotacion)
	return string(evaluateResult),nil
}	

func newVoteHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
    	UiVotacion string `json:"uiVotacion"`;
		IdVoter string `json:"idVoter"`;
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	println(body.UiVotacion)
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error de conexión con el Gateway: %v", err), http.StatusInternalServerError)
		return
	}
	defer gw.Close()

	// Configuración del contrato
	chaincodeName := "votacion"
	channelName := "marketplace"
	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	// Consultar los votos de un candidato
	result, err := newVoter(contract, body.UiVotacion, body.IdVoter) // Usar el nombre dinámico
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener a los votantes: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"elección": result})
}

func getElectionHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
    	UiVotacion string `json:"uiVotacion"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	println(body.UiVotacion)
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error de conexión con el Gateway: %v", err), http.StatusInternalServerError)
		return
	}
	defer gw.Close()

	// Configuración del contrato
	chaincodeName := "votacion"
	channelName := "marketplace"
	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	// Consultar los votos de un candidato
	result, err := getElection(contract, body.UiVotacion) // Usar el nombre dinámico
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener a los votantes: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"elección": result})
}

func getVotesHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
    	UiVotacion string `json:"uiVotacion"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	println(body.UiVotacion)
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error de conexión con el Gateway: %v", err), http.StatusInternalServerError)
		return
	}
	defer gw.Close()

	// Configuración del contrato
	chaincodeName := "votacion"
	channelName := "marketplace"
	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	// Consultar los votos de un candidato
	result, err := getVotes(contract, body.UiVotacion) // Usar el nombre dinámico
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener a los votantes: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"votantes": result})
}

func getCandidatesHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
    	UiVotacion string `json:"uiVotacion"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	println(body.UiVotacion)
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithHash(hash.SHA256),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error de conexión con el Gateway: %v", err), http.StatusInternalServerError)
		return
	}
	defer gw.Close()

	// Configuración del contrato
	chaincodeName := "votacion"
	channelName := "marketplace"
	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	// Consultar los votos de un candidato
	result, err := getCandidates(contract, body.UiVotacion) // Usar el nombre dinámico
	if err != nil {
		http.Error(w, fmt.Sprintf("Error al obtener los votos: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"candidatos": result})
}

func getElection(contract *client.Contract, uiVotacion string) (interface{}, error) {
	evaluateResult, err := contract.EvaluateTransaction("ObtenerVotacion", uiVotacion)
	if err != nil {
		return nil, err
	}
	println(uiVotacion)
	// Si el resultado es un número, lo convertimos en string o lo retornamos directamente
	return string(evaluateResult), nil // o lo procesamos como int si es un número entero
}

func newVoter(contract *client.Contract, uiVotacion string, idVotante string) (interface{}, error) {
	evaluateResult, err := contract.SubmitTransaction("AgregarVotante", uiVotacion, idVotante)
	if err != nil {
		return nil, err
	}
	println(uiVotacion)
	return string(evaluateResult),nil
}

func getVotes(contract *client.Contract, uiVotacion string) (interface{}, error) {
	evaluateResult, err := contract.EvaluateTransaction("ObtenerVotantes", uiVotacion)
	if err != nil {
		return nil, err
	}
	println(uiVotacion)
	// Si el resultado es un número, lo convertimos en string o lo retornamos directamente
	return string(evaluateResult), nil // o lo procesamos como int si es un número entero
}

func getCandidates(contract *client.Contract, uiVotacion string) (interface{}, error) {
	evaluateResult, err := contract.EvaluateTransaction("obtenerCandidatos", uiVotacion)
	if err != nil {
		return nil, err
	}
	println(uiVotacion)
	// Si el resultado es un número, lo convertimos en string o lo retornamos directamente
	return string(evaluateResult), nil // o lo procesamos como int si es un número entero
}

// Nueva función para crear la conexión gRPC
func newGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certifcate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// Nueva función para crear una identidad X.509
func newIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(fmt.Errorf("failed to parse certificate: %w", err))
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(fmt.Errorf("failed to create X509 identity: %w", err))
	}
	return id
}

// Función para leer el primer archivo desde un directorio dado
func readFirstFile(directory string) ([]byte, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in directory: %s", directory)
	}

	return ioutil.ReadFile(path.Join(directory, files[0].Name()))
}

func newSign() identity.Sign {
	keyFiles, err := ioutil.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read key path directory: %w", err))
	}

	if len(keyFiles) == 0 {
		panic(fmt.Errorf("no key file found in directory: %s", keyPath))
	}

	privateKeyPEM, err := ioutil.ReadFile(path.Join(keyPath, keyFiles[0].Name()))
	if err != nil {
		panic(fmt.Errorf("failed to read private key: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(fmt.Errorf("failed to load private key: %w", err))
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(fmt.Errorf("failed to create signer: %w", err))
	}
	return sign
}