package greeting

import "fmt"

// Service define la interfaz para el servicio de saludos
type Service interface {
	GenerateGreeting(name string) (*GreetingResponse, error)
}

type service struct{}

// NewService crea una nueva instancia del servicio
func NewService() Service {
	return &service{}
}

// GenerateGreeting genera un saludo personalizado
func (s *service) GenerateGreeting(name string) (*GreetingResponse, error) {
	if name == "" {
		return nil, fmt.Errorf("el nombre no puede estar vacío")
	}
	
	message := fmt.Sprintf("Hola %s desde Go!", name)
	
	return &GreetingResponse{
		Message: message,
		Status:  "success",
	}, nil
}