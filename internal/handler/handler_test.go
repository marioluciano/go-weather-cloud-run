package handler

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "Zero Celsius",
			celsius:  0,
			expected: 32,
		},
		{
			name:     "100 Celsius",
			celsius:  100,
			expected: 212,
		},
		{
			name:     "Negative Celsius",
			celsius:  -40,
			expected: -40,
		},
		{
			name:     "Room Temperature",
			celsius:  20,
			expected: 68,
		},
		{
			name:     "Decimal Celsius",
			celsius:  23.5,
			expected: 74.3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := celsiusToFahrenheit(tt.celsius)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("celsiusToFahrenheit(%f) = %f, want %f", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		name     string
		celsius  float64
		expected float64
	}{
		{
			name:     "Zero Celsius",
			celsius:  0,
			expected: 273,
		},
		{
			name:     "100 Celsius",
			celsius:  100,
			expected: 373,
		},
		{
			name:     "Negative Celsius",
			celsius:  -273,
			expected: 0,
		},
		{
			name:     "Room Temperature",
			celsius:  20,
			expected: 293,
		},
		{
			name:     "Decimal Celsius",
			celsius:  25.5,
			expected: 298.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := celsiusToKelvin(tt.celsius)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("celsiusToKelvin(%f) = %f, want %f", tt.celsius, result, tt.expected)
			}
		})
	}
}

func TestIsValidCEP(t *testing.T) {
	tests := []struct {
		name     string
		cep      string
		expected bool
	}{
		{
			name:     "Valid CEP",
			cep:      "01310100",
			expected: true,
		},
		{
			name:     "Valid CEP 2",
			cep:      "20040020",
			expected: true,
		},
		{
			name:     "Invalid CEP - Too short",
			cep:      "0131010",
			expected: false,
		},
		{
			name:     "Invalid CEP - Too long",
			cep:      "013101000",
			expected: false,
		},
		{
			name:     "Invalid CEP - Contains letters",
			cep:      "0131010A",
			expected: false,
		},
		{
			name:     "Invalid CEP - Contains symbols",
			cep:      "01310-100",
			expected: false,
		},
		{
			name:     "Empty CEP",
			cep:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidCEP(tt.cep)
			if result != tt.expected {
				t.Errorf("isValidCEP(%s) = %v, want %v", tt.cep, result, tt.expected)
			}
		})
	}
}

// MockViaCEPClient is a mock implementation of the ViaCEP client
type MockViaCEPClient struct {
	MockFetchByCep func(ctx any, cep string) (any, error)
}

// MockWeatherAPIClient is a mock implementation of the WeatherAPI client
type MockWeatherAPIClient struct {
	MockFetchByCity func(ctx any, city string) (any, error)
}

func TestHandleWeatherRequest_MethodNotAllowed(t *testing.T) {
	handler := &WeatherHandler{}

	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	rec := httptest.NewRecorder()

	handler.HandleWeatherRequest(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}

	var errResp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if errResp.Message != "method not allowed" {
		t.Errorf("Expected error message 'method not allowed', got '%s'", errResp.Message)
	}
}

func TestHandleWeatherRequest_InvalidRequestFormat(t *testing.T) {
	handler := &WeatherHandler{}

	body := bytes.NewBufferString("invalid json")
	req := httptest.NewRequest(http.MethodPost, "/weather", body)
	rec := httptest.NewRecorder()

	handler.HandleWeatherRequest(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var errResp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if errResp.Message != "invalid request format" {
		t.Errorf("Expected error message 'invalid request format', got '%s'", errResp.Message)
	}
}

func TestHandleWeatherRequest_InvalidZipCode(t *testing.T) {
	handler := &WeatherHandler{}

	reqBody := WeatherRequest{CEP: "12345"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/weather", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.HandleWeatherRequest(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, rec.Code)
	}

	var errResp ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if errResp.Message != "invalid zipcode" {
		t.Errorf("Expected error message 'invalid zipcode', got '%s'", errResp.Message)
	}
}

func TestWeatherResponse_Structure(t *testing.T) {
	resp := WeatherResponse{
		TempC: 20.0,
		TempF: 68.0,
		TempK: 293.0,
	}

	// Verify the structure
	if resp.TempC != 20.0 {
		t.Errorf("Expected TempC 20.0, got %f", resp.TempC)
	}
	if resp.TempF != 68.0 {
		t.Errorf("Expected TempF 68.0, got %f", resp.TempF)
	}
	if resp.TempK != 293.0 {
		t.Errorf("Expected TempK 293.0, got %f", resp.TempK)
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var unmarshaled WeatherResponse
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if unmarshaled.TempC != resp.TempC || unmarshaled.TempF != resp.TempF || unmarshaled.TempK != resp.TempK {
		t.Errorf("JSON marshaling/unmarshaling failed")
	}
}

func TestWeatherRequest_Parsing(t *testing.T) {
	jsonStr := `{"cep":"01310100"}`

	var req WeatherRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Failed to unmarshal request: %v", err)
	}

	if req.CEP != "01310100" {
		t.Errorf("Expected CEP '01310100', got '%s'", req.CEP)
	}
}

func TestTemperatureRounding(t *testing.T) {
	tests := []struct {
		name    string
		celsius float64
		expectC float64
		expectF float64
		expectK float64
	}{
		{
			name:    "Single decimal place",
			celsius: 25.123,
			expectC: 25.1,
			expectF: 77.2,
			expectK: 298.1,
		},
		{
			name:    "Rounding up",
			celsius: 19.95,
			expectC: 20.0,
			expectF: 67.9,
			expectK: 293.0,
		},
		{
			name:    "Exact values",
			celsius: 25,
			expectC: 25.0,
			expectF: 77.0,
			expectK: 298.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempC := math.Round(tt.celsius*10) / 10
			tempF := math.Round((tt.celsius*1.8+32)*10) / 10
			tempK := math.Round((tt.celsius+273)*10) / 10

			if math.Abs(tempC-tt.expectC) > 0.01 {
				t.Errorf("TempC: expected %f, got %f", tt.expectC, tempC)
			}
			if math.Abs(tempF-tt.expectF) > 0.01 {
				t.Errorf("TempF: expected %f, got %f", tt.expectF, tempF)
			}
			if math.Abs(tempK-tt.expectK) > 0.01 {
				t.Errorf("TempK: expected %f, got %f", tt.expectK, tempK)
			}
		})
	}
}
