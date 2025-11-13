package greeting

// GreetingRequest representa la solicitud para saludar
type GreetingRequest struct {
	Name string `json:"name"`
}

// GreetingResponse representa la respuesta del saludo
type GreetingResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}