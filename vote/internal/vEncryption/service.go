package vEncryption

// import (
// 	"bytes"
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"crypto/rand"
// 	"crypto/sha256"
// 	"encoding/json"
// 	"fmt"
// 	"goHEncryption/pkg/config"
// 	"io"
// 	"log"
// 	"mime/multipart"
// 	"net/http"
// 	"os"

// 	"github.com/tuneinsight/lattigo/v6/core/rlwe"
// 	"github.com/tuneinsight/lattigo/v6/multiparty"
// 	"github.com/tuneinsight/lattigo/v6/schemes/bgv"
// 	"github.com/tuneinsight/lattigo/v6/utils/sampling"
// )

// // Service define la interfaz para el servicio de saludos
// type Service interface {
// 	GenerateGreeting(name string) (*GreetingResponse, error)
// 	Nuevo() (*GreetingResponse, error)
// 	GenerateKeys(Init Init) (*GreetingResponse, error)
// }

// type service struct{}

// // NewService crea una nueva instancia del servicio
// func NewService() Service {
// 	return &service{}
// }

// // Prueba
// func (s *service) Nuevo() (*GreetingResponse, error) {

// 	message := fmt.Sprintln("Hola desde Go!")

// 	return &GreetingResponse{
// 		Message: message,
// 		Status:  "success",
// 	}, nil
// }

// // GenerateGreeting genera un saludo personalizado
// func (s *service) GenerateGreeting(name string) (*GreetingResponse, error) {
// 	if name == "" {
// 		return nil, fmt.Errorf("el nombre no puede estar vacío")
// 	}

// 	message := fmt.Sprintf("Hola %s desde Go!", name)

// 	return &GreetingResponse{
// 		Message: message,
// 		Status:  "success",
// 	}, nil
// }


// func encryptAESGCM(data []byte, key []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	aesGCM, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, err
// 	}
// 	nonce := make([]byte, aesGCM.NonceSize())
// 	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
// 		return nil, err
// 	}
// 	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)
// 	return ciphertext, nil
// }


// func sendAE(ae Ae, party Party) error {
// url := "http://127.0.0.1:8000/send/ae"

// fmt.Printf("%+v\n", party)
// 	// === 1. Serializar la party a JSON ===
// 	partyData, err := json.MarshalIndent(party, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("error serializando party: %v", err)
// 	}

// 	fmt.Println(string(partyData))
// 	// === 2. Encriptar el archivo con el hash del password ===
// 	password := "123"
// 	key := sha256.Sum256([]byte(password)) // Derivar la clave AES desde el password
// 	cipherData, err := encryptAESGCM(partyData, key[:])
// 	if err != nil {
// 		return fmt.Errorf("error cifrando archivo: %v", err)
// 	}

// 	// === 3. Guardar archivo temporal ===
// 	tempFile, err := os.CreateTemp("", "party_*.enc")
// 	if err != nil {
// 		return fmt.Errorf("error creando archivo temporal: %v", err)
// 	}
// 	defer os.Remove(tempFile.Name()) // eliminar después de enviar

// 	if _, err := tempFile.Write(cipherData); err != nil {
// 		return fmt.Errorf("error escribiendo archivo cifrado: %v", err)
// 	}
// 	tempFile.Close()

// 	// === 4. Construir formulario ===
// 	body := &bytes.Buffer{}
// 	writer := multipart.NewWriter(body)

// 	writer.WriteField("nombre", ae.Name)
// 	writer.WriteField("destinatario", ae.Mail)
// 	writer.WriteField("password", password)

// 	file, err := os.Open(tempFile.Name())
// 	if err != nil {
// 		return fmt.Errorf("error abriendo archivo cifrado: %v", err)
// 	}
// 	defer file.Close()

// 	part, err := writer.CreateFormFile("archivo", "party_encrypted.enc")
// 	if err != nil {
// 		return fmt.Errorf("error creando parte de archivo: %v", err)
// 	}
// 	if _, err := io.Copy(part, file); err != nil {
// 		return fmt.Errorf("error copiando archivo al form: %v", err)
// 	}

// 	writer.Close()

// 	// === 5. Enviar la solicitud ===
// 	req, err := http.NewRequest("POST", url, body)
// 	if err != nil {
// 		return fmt.Errorf("error creando solicitud: %v", err)
// 	}
// 	req.Header.Set("Content-Type", writer.FormDataContentType())

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("error enviando solicitud: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("error en respuesta: %s", resp.Status)
// 	}

// 	log.Printf("✅ AE enviada correctamente a %s", ae.Mail)
// 	return nil
// }


// func (s *service) GenerateKeys(init Init) (*GreetingResponse, error) {
// 	list := init.N
// 	t := init.T
// 	parties := createParties(len(list), t)

// 	for i := 0; i < len(parties); i++ {
// 		list[i].IdParty = parties[i].Id
// 		//enviar parties por esta direcion http://127.0.0.1:8000/send/ae
		
// 		sendAE(list[i], *parties[i])


// 	}

// 	// guardar Init en base de datos

// 	message := fmt.Sprintln("Procesado")

// 	return &GreetingResponse{
// 		Message: message,
// 		Status:  "success",
// 	}, nil

// }


// // ======================== PUBLIC KEY GENERATION ========================

// func collectivePublicKey(parties []*Party, t int, params bgv.Parameters) *rlwe.PublicKey {
// 	crs, err := sampling.NewPRNG()
// 	if err != nil {
// 		log.Fatalf("error creando PRNG: %v", err)
// 	}

// 	ckg := multiparty.NewPublicKeyGenProtocol(params)
// 	crp := ckg.SampleCRP(crs)

// 	activeParties := parties[:t]
// 	shamirPts := getShamirPoints(activeParties)

// 	shares := make([]multiparty.PublicKeyGenShare, len(activeParties))
// 	tsks := make([]*rlwe.SecretKey, len(activeParties))

// 	for i, p := range activeParties {
// 		shares[i] = ckg.AllocateShare()
// 		tsks[i] = rlwe.NewSecretKey(params)
		
// 		err := p.Combiner.GenAdditiveShare(shamirPts, p.ShamirPt, p.Tsk, tsks[i])
// 		if err != nil {
// 			log.Fatal(err)
// 		}
		
// 		ckg.GenShare(tsks[i], crp, &shares[i])
// 	}

// 	aggShare := shares[0]
// 	for i := 1; i < len(activeParties); i++ {
// 		ckg.AggregateShares(aggShare, shares[i], &aggShare)
// 	}

// 	publicKey := rlwe.NewPublicKey(params)
// 	ckg.GenPublicKey(aggShare, crp, publicKey)

// 	//guardar publickey

// 	return publicKey
// }

// func getShamirPoints(parties []*Party) []multiparty.ShamirPublicPoint {
// 	pts := make([]multiparty.ShamirPublicPoint, len(parties))
// 	for i, p := range parties {
// 		pts[i] = p.ShamirPt
// 	}
// 	return pts
// }

// func createParties(N, t int) []*Party {
// 	SaveParams()

// 	params, err := LoadParams()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	gen := rlwe.NewKeyGenerator(params)
// 	parties := make([]*Party, N)
// 	shamirPts := make([]multiparty.ShamirPublicPoint, N)

// 	for i := 0; i < N; i++ {
// 		parties[i] = &Party{
// 			Id:       i,
// 			Sk:       gen.GenSecretKeyNew(),
// 			ShamirPt: multiparty.ShamirPublicPoint(i + 1),
// 		}
// 		shamirPts[i] = parties[i].ShamirPt
// 	}

// 	thresholdizer := multiparty.NewThresholdizer(params)

// 	for i := range parties {
// 		poly, err := thresholdizer.GenShamirPolynomial(t, parties[i].Sk)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		for j := range parties {
// 			share := thresholdizer.AllocateThresholdSecretShare()
// 			thresholdizer.GenShamirSecretShare(parties[j].ShamirPt, poly, &share)

// 			if i == 0 {
// 				parties[j].Tsk = share
// 			} else {
// 				thresholdizer.AggregateShares(parties[j].Tsk, share, &parties[j].Tsk)
// 			}
// 		}
// 	}

// 	for i := range parties {
// 		parties[i].Combiner = multiparty.NewCombiner(*params.GetRLWEParameters(), parties[i].ShamirPt, shamirPts, t)
// 	}

// 	// Enviar todo a las AE

// 	return parties
// }

// func initParams() bgv.Parameters {
// 	// Cargar cofiguracion
// 	cfg := config.Load()

// 	// Usamos un módulo primo grande que funcione con BGV
// 	params, err := bgv.NewParametersFromLiteral(bgv.ParametersLiteral{
// 		LogN:             cfg.He.LogN,
// 		LogQ:             []int{cfg.He.LogQ, cfg.He.LogQ, cfg.He.LogQ, cfg.He.LogQ},
// 		LogP:             []int{cfg.He.LogP},
// 		PlaintextModulus: uint64(cfg.He.Module),
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return params
// }

// func SaveParams() error {
// 	params := initParams()
// 	filename := "params.json"
// 	data, err := json.Marshal(params.ParametersLiteral())
// 	if err != nil {
// 		return err
// 	}
// 	return os.WriteFile(filename, data, 0600)
// }

// func LoadParams() (bgv.Parameters, error) {
// 	filename := "params.json"
// 	data, err := os.ReadFile(filename)
// 	if err != nil {
// 		return bgv.Parameters{}, err
// 	}

// 	var literal bgv.ParametersLiteral
// 	if err := json.Unmarshal(data, &literal); err != nil {
// 		return bgv.Parameters{}, err
// 	}

// 	return bgv.NewParametersFromLiteral(literal)
// }



// // ======================== ENCRYPTION ========================

// func encryptMessage(message uint64, pk *rlwe.PublicKey, params bgv.Parameters) *rlwe.Ciphertext {
// 	encoder := bgv.NewEncoder(params)
// 	encryptor := rlwe.NewEncryptor(params, pk)
	
// 	// Verificar que el mensaje está dentro del módulo
// 	modulus := uint64(params.PlaintextModulus())
// 	if message >= modulus {
// 		log.Fatalf("❌ Número %d excede el módulo %d", message, modulus)
// 	}
	
// 	// Codificar el número en todos los slots
// 	vector := make([]uint64, params.N())
// 	for i := range vector {
// 		vector[i] = message
// 	}
	
// 	pt := bgv.NewPlaintext(params, params.MaxLevel())
// 	encoder.Encode(vector, pt)
	
// 	ct, err := encryptor.EncryptNew(pt)
// 	if err != nil {
// 		log.Fatalf("error cifrando el mensaje: %v", err)
// 	}
// 	return ct
// }

// // ======================== COLLECTIVE DECRYPTION ========================

// func decryptCollective(ct *rlwe.Ciphertext, parties []*Party, t int, params bgv.Parameters) uint64 {
// 	if len(parties) < t {
// 		log.Fatalf("No hay suficientes partes (%d) para el umbral t=%d", len(parties), t)
// 	}

// 	collectiveSK := rlwe.NewSecretKey(params)
// 	shamirPts := getShamirPoints(parties)

// 	for _, p := range parties {
// 		skShare := rlwe.NewSecretKey(params)
// 		err := p.Combiner.GenAdditiveShare(shamirPts, p.ShamirPt, p.Tsk, skShare)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		params.RingQP().Add(collectiveSK.Value, skShare.Value, collectiveSK.Value)
// 	}

// 	decryptor := rlwe.NewDecryptor(params, collectiveSK)
// 	pt := bgv.NewPlaintext(params, params.MaxLevel())
// 	decryptor.Decrypt(ct, pt)

// 	encoder := bgv.NewEncoder(params)
// 	decoded := make([]uint64, params.N())
// 	encoder.Decode(pt, decoded)
	
// 	return decoded[0]
// }

// // ======================== OPERACIONES SOBRE CIPHERTEXTS ========================

// func sumarCiphertexts(params bgv.Parameters, ciphertexts []*rlwe.Ciphertext) *rlwe.Ciphertext {
// 	if len(ciphertexts) == 0 {
// 		return nil
// 	}
// 	evaluator := bgv.NewEvaluator(params, nil)
// 	resultado := ciphertexts[0].CopyNew()
// 	for i := 1; i < len(ciphertexts); i++ {
// 		if err := evaluator.Add(resultado, ciphertexts[i], resultado); err != nil {
// 			log.Fatal("Error sumando ciphertexts:", err)
// 		}
// 	}
// 	return resultado
// }

