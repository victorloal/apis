package auth

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

func (s *Service) GenerateToken(serviceName string) (string, error) {
	claims := jwt.MapClaims{
		"service": serviceName,
		"iss":     "auth-service",
		"exp":     time.Now().Add(5 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) ValidateHMAC(message []byte, signature string) bool {
	mac := hmac.New(sha256.New, s.hmacSecret)
	mac.Write(message)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}