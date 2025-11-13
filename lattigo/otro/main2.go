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

// ------------------ Setup ------------------

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

func getShamirPoints(parties []*Party) []multiparty.ShamirPublicPoint {
	pts := make([]multiparty.ShamirPublicPoint, len(parties))
	for i, p := range parties {
		pts[i] = p.shamirPt
	}
	return pts
}

// ------------------ Generar PK colectiva ------------------

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

// ------------------ Cifrado y suma ------------------

func encryptMessageUnderCollective(message uint64, pk *rlwe.PublicKey, params bgv.Parameters) *rlwe.Ciphertext {
	encoder := bgv.NewEncoder(params)
	vector := make([]uint64, params.N())
	for i := range vector {
		vector[i] = message
	}
	pt := bgv.NewPlaintext(params, params.MaxLevel())
	encoder.Encode(vector, pt)

	encryptor := rlwe.NewEncryptor(params, pk)
	ct, err := encryptor.EncryptNew(pt)
	if err != nil {
		log.Fatalf("encrypt failed: %v", err)
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

// Reencriptado con un solo ciphertext
func reencryptWithSingleCiphertext(ctSum, ctOneSwitched *rlwe.Ciphertext, params bgv.Parameters) *rlwe.Ciphertext {
    evaluator := bgv.NewEvaluator(params, nil)
    result := rlwe.NewCiphertext(params, ctSum.Degree(), ctSum.Level())
    if err := evaluator.Mul(ctSum, ctOneSwitched, result); err != nil {
        log.Fatalf("error applying single-ciphertext reencryption: %v", err)
    }
    evaluator.Relinearize(result, result) // opcional si tu esquema requiere
    return result
}

// ------------------ MAIN ------------------

func main() {
	params := initParams()
N := 12
t := 10

fmt.Println("Iniciando...")

// 1) Crear parties y PK colectiva
parties := createParties(N, t, params)
publicKey := collectivePublicKey(parties, t, params)
fmt.Println("PK_collective generada")

// 2) Generar PK_dest / SK_dest (verificador único)
kg := rlwe.NewKeyGenerator(params)
skDest := kg.GenSecretKeyNew()
pkDest := kg.GenPublicKeyNew(skDest)
fmt.Println("PK_dest generada")

// 3) Crear "ciphertext de 1" bajo PK_dest (plantilla para reencriptado)
encoder := bgv.NewEncoder(params)
ptOne := bgv.NewPlaintext(params, params.MaxLevel())
vec := make([]uint64, params.N())
for i := range vec {
    vec[i] = 1
}
encoder.Encode(vec, ptOne)

encryptor := rlwe.NewEncryptor(params, pkDest)
ctOneSwitched, err := encryptor.EncryptNew(ptOne)
if err != nil {
    log.Fatalf("encrypt ctOneSwitched failed: %v", err)
}
fmt.Println("Ciphertext de 1 creado bajo PK_dest")

// 4) Cifrar mensajes bajo PK_colectiva
msgs := []uint64{0, 0, 0, 1, 0}
ct := make([]*rlwe.Ciphertext, len(msgs))
for i, m := range msgs {
    ct[i] = encryptMessageUnderCollective(m, publicKey, params)
}
fmt.Println("Mensajes cifrados bajo PK_collective")

// 5) Sumar homomórficamente
ctSum := sumarCiphertexts(params, ct)
fmt.Println("ctSum listo")

// 6) Reencriptado usando el ciphertext de 1
evaluator := bgv.NewEvaluator(params, nil)
ctSwitched := rlwe.NewCiphertext(params, ctSum.Degree(), ctSum.Level())
if err := evaluator.Mul(ctSum, ctOneSwitched, ctSwitched); err != nil {
    log.Fatalf("error aplicando reencriptado: %v", err)
}
fmt.Println("Reencriptado aplicado -> ctSwitched bajo PK_dest")

// 7) Desencriptar resultado final con SK_dest
decryptor := rlwe.NewDecryptor(params, skDest)
pt := bgv.NewPlaintext(params, ctSwitched.Level())
decryptor.Decrypt(ctSwitched, pt)

decoded := make([]uint64, params.N())
encoder.Decode(pt, decoded)

fmt.Println("Resultado (slot0):", decoded[0])
if decoded[0] == 1 {
    fmt.Println("VERIFICACIÓN OK")
} else {
    log.Fatalf("ERROR: suma != 1 (valor=%d)", decoded[0])
}
}
