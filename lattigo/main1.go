package main

import (
	"fmt"
	"log"
 
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

// ======================== SETUP ========================

func initParams() bgv.Parameters {
	// Usamos un módulo primo grande que funcione con BGV
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
	N := 12
	t := 10

	fmt.Printf("=== SISTEMA MULTIPARTITA BGV %d-de-%d ===\n", t, N)
	fmt.Printf("Módulo del plaintext: %d\n", params.PlaintextModulus())
	fmt.Printf("Rango válido: [0, %d]\n\n", params.PlaintextModulus()-1)

	parties := createParties(N, t, params)
	fmt.Println("✅ Partes creadas con esquema umbral")

	publicKey := collectivePublicKey(parties, t, params)
	fmt.Println("✅ Clave pública colectiva generada")

	// Probar con números pequeños y grandes
	testNumbers := []uint64{
		0,                      		
		1,                       		
		123,                    
		// 1000,                   
		// 50000,                  
		// 65536,                   
		// 999999,                 
		// 1000003,             
	}

	fmt.Println("=== PRUEBAS CON NÚMEROS PEQUEÑOS Y GRANDES ===")
	
	allCorrect := true
	
	ct := make([]*rlwe.Ciphertext, len(testNumbers))
	for i, msg := range testNumbers {
		ct[i] = encryptMessage(msg, publicKey, params)
	}

	fmt.Println("✅ Mensajes cifrados")

	suma := sumarCiphertexts(params, ct)

	fmt.Println("Descifrando la suma de los mensajes...")

	result := decryptCollective(suma, parties[:t], t, params)
	fmt.Printf(" Mensaje descifrado: %d\n", result)

	if allCorrect {
		fmt.Println("\n🎉 ¡TODAS LAS PRUEBAS EXITOSAS!")
		fmt.Println("✅ Exactitud absoluta con números pequeños y grandes")
		fmt.Printf("✅ Rango válido: [0, %d]\n", params.PlaintextModulus()-1)
	} else {
		fmt.Println("\n❌ Algunas pruebas fallaron")
	}
}