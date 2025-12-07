package ledger

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	jwtSecret  []byte
	hmacSecret []byte
}

func NewService(jwtSecret, hmacSecret string) *Service {
	return &Service{
		jwtSecret:  []byte(jwtSecret),
		hmacSecret: []byte(hmacSecret),
	}
}

func (s *Service) ValidateJWT(tokenStr string) bool {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	return err == nil && token.Valid
}

func (s *Service) ValidateHMAC(message []byte, signature string) bool {
	mac := hmac.New(sha256.New, s.hmacSecret)
	mac.Write(message)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}

func (s *Service) ProcessTransaction(tx Transaction) TransactionResponse {
	return TransactionResponse{
		Status:    "confirmed",
		TxID:      "tx_" + time.Now().Format("20060102150405"),
		Timestamp: time.Now(),
		Ledger:    "secure-blockchain",
	}
}