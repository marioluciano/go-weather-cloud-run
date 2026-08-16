package viacep

type ViaCEPResponse struct {
	City  string `json:"localidade"`
	Error string `json:"erro,omitempty"`
}
