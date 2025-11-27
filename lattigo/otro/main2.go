package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/multiparty"
	"github.com/tuneinsight/lattigo/v6/schemes/bgv"
	"github.com/tuneinsight/lattigo/v6/utils/sampling"
)

type Party struct {
	id       int
	sk       *rlwe.SecretKey
	shamirPt multiparty.ShamirPublicPoint
	tsk      multiparty.ShamirSecretShare
	combiner multiparty.Combiner
}

// ======================== SERIALIZACIÓN ========================

// CiphertextData representa un ciphertext serializado
type CiphertextData struct {
	EncryptedData string `json:"encrypted_data"` // Base64 o string del ciphertext
	OriginalValue uint64 `json:"original_value,omitempty"`
}

// EncryptedCollection representa la colección completa cifrada
type EncryptedCollection struct {
	Parameters struct {
		LogN             int   `json:"log_n"`
		PlaintextModulus uint64 `json:"plaintext_modulus"`
	} `json:"parameters"`
	Ciphertexts []CiphertextData `json:"ciphertexts"`
	TotalCount  int              `json:"total_count"`
}

// ciphertextToString convierte un ciphertext a string (simplificado)
func ciphertextToString(ct *rlwe.Ciphertext) string {
	// Para una representación simple, usamos MarshalBinary
	data, err := ct.MarshalBinary()
	if err != nil {
		return ""
	}
	
	// Convertir a string (en producción usarías base64)
	// Esta es una versión simplificada para demostración
	return fmt.Sprintf("%x", data)
}

// stringToCiphertext convierte string a ciphertext
func stringToCiphertext(dataStr string, params bgv.Parameters) (*rlwe.Ciphertext, error) {
	ct := rlwe.NewCiphertext(params, 1, params.MaxLevel())
	
	// Convertir de string hexadecimal a bytes
	// En producción usarías base64
	var data []byte
	_, err := fmt.Sscanf(dataStr, "%x", &data)
	if err != nil {
		return nil, err
	}
	
	err = ct.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}
	
	return ct, nil
}

// saveEncryptedToJSON guarda los ciphertexts en formato array de strings
func saveEncryptedToJSON(ciphertexts []*rlwe.Ciphertext, filename string) error {
	// Crear array de strings con los ciphertexts
	encryptedStrings := make([]string, len(ciphertexts))
	for i, ct := range ciphertexts {
		encryptedStrings[i] = ciphertextToString(ct)
	}

	// Convertir a JSON array
	jsonData, err := json.Marshal(encryptedStrings)
	if err != nil {
		return fmt.Errorf("error codificando JSON: %w", err)
	}

	// Guardar en archivo
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creando archivo: %w", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error escribiendo archivo: %w", err)
	}

	return nil
}


// loadEncryptedFromJSON carga los ciphertexts desde JSON
func loadEncryptedFromJSON(filename string, params bgv.Parameters) ([]*rlwe.Ciphertext, []uint64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error abriendo archivo: %w", err)
	}
	defer file.Close()

	var collection EncryptedCollection
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&collection); err != nil {
		return nil, nil, fmt.Errorf("error decodificando JSON: %w", err)
	}

	ciphertexts := make([]*rlwe.Ciphertext, len(collection.Ciphertexts))
	originalNumbers := make([]uint64, len(collection.Ciphertexts))

	for i, ctData := range collection.Ciphertexts {
		ct, err := stringToCiphertext(ctData.EncryptedData, params)
		if err != nil {
			return nil, nil, fmt.Errorf("error convirtiendo ciphertext %d: %w", i, err)
		}
		ciphertexts[i] = ct
		originalNumbers[i] = ctData.OriginalValue
	}

	return ciphertexts, originalNumbers, nil
}

// ======================== SETUP ========================

func initParams() bgv.Parameters {
	// Usamos un módulo primo grande que funcione con BGV
	params, err := bgv.NewParametersFromLiteral(bgv.ParametersLiteral{
		LogN:             15, // Reducido para mejor rendimiento
		LogQ:             []int{54, 54, 54},
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
				parties[j].tsk = share
			} else {
				thresholdizer.AggregateShares(parties[j].tsk, share, &parties[j].tsk)
			}
		}
	}

	for i := range parties {
		parties[i].combiner = multiparty.NewCombiner(*params.GetRLWEParameters(), parties[i].shamirPt, shamirPts, t)
	}

	return parties
}

// ======================== PUBLIC KEY GENERATION ========================

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
		
		err := p.combiner.GenAdditiveShare(shamirPts, p.shamirPt, p.tsk, tsks[i])
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

func getShamirPoints(parties []*Party) []multiparty.ShamirPublicPoint {
	pts := make([]multiparty.ShamirPublicPoint, len(parties))
	for i, p := range parties {
		pts[i] = p.shamirPt
	}
	return pts
}

// ======================== ENCRYPTION ========================

func encryptMessage(message uint64, pk *rlwe.PublicKey, params bgv.Parameters) *rlwe.Ciphertext {
	encoder := bgv.NewEncoder(params)
	encryptor := rlwe.NewEncryptor(params, pk)
	
	// Verificar que el mensaje está dentro del módulo
	modulus := uint64(params.PlaintextModulus())
	if message >= modulus {
		log.Fatalf("❌ Número %d excede el módulo %d", message, modulus)
	}
	
	// Codificar el número en todos los slots
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

// ======================== COLLECTIVE DECRYPTION ========================

func decryptCollective(ct *rlwe.Ciphertext, parties []*Party, t int, params bgv.Parameters) uint64 {
	if len(parties) < t {
		log.Fatalf("No hay suficientes partes (%d) para el umbral t=%d", len(parties), t)
	}

	collectiveSK := rlwe.NewSecretKey(params)
	shamirPts := getShamirPoints(parties)

	for _, p := range parties {
		skShare := rlwe.NewSecretKey(params)
		err := p.combiner.GenAdditiveShare(shamirPts, p.shamirPt, p.tsk, skShare)
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

// ======================== OPERACIONES SOBRE CIPHERTEXTS ========================

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

// ======================== MAIN FLOW ========================

func main() {
	params := initParams()
	N := 6  // Reducido para mejor rendimiento
	t := 4  // Reducido para mejor rendimiento

	fmt.Printf("=== SISTEMA MULTIPARTITA BGV %d-de-%d ===\n", t, N)
	fmt.Printf("Módulo del plaintext: %d\n", params.PlaintextModulus())
	fmt.Printf("Rango válido: [0, %d]\n\n", params.PlaintextModulus()-1)

	parties := createParties(N, t, params)
	fmt.Println("✅ Partes creadas con esquema umbral")

	publicKey := collectivePublicKey(parties, t, params)
	fmt.Println("✅ Clave pública colectiva generada")

	// Números a cifrar
	testNumbers := []uint64{0, 0, 1, 0}

	fmt.Printf("=== CIFRANDO NÚMEROS: %v ===\n", testNumbers)
	
	// Cifrar los números
	encryptedCiphertexts := make([]*rlwe.Ciphertext, 0)
	for _, msg := range testNumbers {
		ct := encryptMessage(msg, publicKey, params)
		encryptedCiphertexts = append(encryptedCiphertexts, ct)
		fmt.Printf("✅ Número %d cifrado correctamente\n", msg)
	}
	
	// Guardar en archivo JSON
	filename := "encrypted_numbers.json"
	err := saveEncryptedToJSON(encryptedCiphertexts, filename)
	if err != nil {
		log.Fatalf("❌ Error guardando archivo: %v", err)
	}
	fmt.Printf("✅ Archivo guardado: %s\n", filename)

	// Leer el archivo JSON
	fmt.Println("\n=== LEYENDO ARCHIVO CIFRADO ===")
	loadedCiphertexts, originalNumbers, err := loadEncryptedFromJSON(filename, params)
	if err != nil {
		log.Fatalf("❌ Error leyendo archivo: %v", err)
	}
	fmt.Printf("✅ Archivo leído: %d ciphertexts cargados\n", len(loadedCiphertexts))

	// Descifrar individualmente
	fmt.Println("\n=== DESCIFRANDO NÚMEROS INDIVIDUALES ===")
	allCorrect := true
	for i, ct := range loadedCiphertexts {
		decrypted := decryptCollective(ct, parties[:t], t, params)
		original := originalNumbers[i]
		
		fmt.Printf("Ciphertext %d: Original=%d, Descifrado=%d", i+1, original, decrypted)
		if decrypted == original {
			fmt.Println(" ✅")
		} else {
			fmt.Println(" ❌")
			allCorrect = false
		}
	}

	// Sumar todos los ciphertexts y descifrar el resultado
	fmt.Println("\n=== SUMANDO TODOS LOS CIPHERTEXTS ===")
	suma := sumarCiphertexts(params, loadedCiphertexts)
	
	// Calcular suma esperada
	sumaEsperada := uint64(0)
	for _, num := range testNumbers {
		sumaEsperada += num
	}

	fmt.Printf("Suma esperada: %d\n", sumaEsperada)
	
	// Descifrar la suma
	resultadoSuma := decryptCollective(suma, parties[:t], t, params)
	fmt.Printf("Suma descifrada: %d", resultadoSuma)
	
	if resultadoSuma == sumaEsperada {
		fmt.Println(" ✅")
	} else {
		fmt.Println(" ❌")
		allCorrect = false
	}

	// Resultado final
	fmt.Println("\n" + strings.Repeat("=", 50))
	if allCorrect {
		fmt.Println("🎉 ¡TODAS LAS PRUEBAS EXITOSAS!")
		fmt.Println("✅ Cifrado, guardado, lectura y descifrado correctos")
		fmt.Printf("✅ Rango válido: [0, %d]\n", params.PlaintextModulus()-1)
	} else {
		fmt.Println("❌ Algunas pruebas fallaron")
	}
	
	// Mostrar contenido del archivo
	fmt.Println("\n=== CONTENIDO DEL ARCHIVO JSON ===")
	fileContent, _ := os.ReadFile(filename)
	fmt.Println(string(fileContent))
}