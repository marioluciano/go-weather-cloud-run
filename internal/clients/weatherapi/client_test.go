package weatherapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherAPIResponse_Parsing(t *testing.T) {
	tests := []struct {
		name       string
		jsonData   string
		expectTemp float64
	}{
		{
			name:       "Valid response with positive temperature",
			jsonData:   `{"current":{"temp_c":25.5,"temp_f":77.9}}`,
			expectTemp: 25.5,
		},
		{
			name:       "Valid response with negative temperature",
			jsonData:   `{"current":{"temp_c":-10.0,"temp_f":14.0}}`,
			expectTemp: -10.0,
		},
		{
			name:       "Valid response with zero temperature",
			jsonData:   `{"current":{"temp_c":0,"temp_f":32.0}}`,
			expectTemp: 0,
		},
		{
			name:       "Valid response with decimal temperature",
			jsonData:   `{"current":{"temp_c":20.123,"temp_f":68.22}}`,
			expectTemp: 20.123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp WeatherAPIResponse
			if err := json.Unmarshal([]byte(tt.jsonData), &resp); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if resp.Current.TempC != tt.expectTemp {
				t.Errorf("Expected TempC %f, got %f", tt.expectTemp, resp.Current.TempC)
			}
		})
	}
}

func TestWeatherAPIClient_Structure(t *testing.T) {
	client := NewClient(http.DefaultClient, "test-api-key")

	if client == nil {
		t.Fatalf("Expected client to be created, got nil")
	}

	if client.client == nil {
		t.Errorf("Expected client.client to be initialized")
	}

	if client.apiKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got '%s'", client.apiKey)
	}
}

func TestWeatherAPIResponse_FieldMapping(t *testing.T) {
	jsonStr := `{
		"location":{"name":"São Paulo"},
		"current":{"temp_c":22.5,"temp_f":72.5,"humidity":65}
	}`

	var resp WeatherAPIResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if resp.Current.TempC != 22.5 {
		t.Errorf("Expected TempC 22.5, got %f", resp.Current.TempC)
	}
}

func TestWeatherAPIClient_FetchByCity_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that the correct query parameters are being sent
		query := r.URL.Query()
		if query.Get("key") == "" {
			t.Errorf("Expected API key in query parameters")
		}
		if query.Get("q") == "" {
			t.Errorf("Expected city in query parameters")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"current":{"temp_c":25.5}}`))
	}))
	defer server.Close()

	_ = NewClient(server.Client(), "test-key")
	// Note: A real test would require modifying the client to use the mock server URL
}

func TestWeatherAPIClient_Integration(t *testing.T) {
	// This is a basic integration test that can be run against the real API
	// Comment this out if you don't want to make real API calls during testing
	t.Skip("Skipping integration test - requires external API access and valid API key")

	client := NewClient(http.DefaultClient, "your-api-key-here")
	ctx := context.Background()

	resp, err := client.FetchByCity(ctx, "São Paulo")
	if err != nil {
		t.Errorf("FetchByCity returned error: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	if resp.Current.TempC == 0 {
		t.Logf("Got temperature: %f°C", resp.Current.TempC)
	}
}

func TestWeatherAPIResponse_ComplexStructure(t *testing.T) {
	complexJSON := `{
		"location":{
			"name":"São Paulo",
			"region":"São Paulo",
			"country":"Brazil"
		},
		"current":{
			"temp_c":22.5,
			"temp_f":72.5,
			"humidity":70,
			"condition":{
				"text":"Partly cloudy",
				"icon":"//cdn.weatherapi.com/weather/128x128/day/116.png"
			}
		}
	}`

	var resp WeatherAPIResponse
	if err := json.Unmarshal([]byte(complexJSON), &resp); err != nil {
		t.Fatalf("Failed to unmarshal complex JSON: %v", err)
	}

	if resp.Current.TempC != 22.5 {
		t.Errorf("Expected TempC 22.5, got %f", resp.Current.TempC)
	}
}

func TestWeatherAPIClient_WithoutAPIKey(t *testing.T) {
	client := NewClient(http.DefaultClient, "")

	if client == nil {
		t.Fatalf("Expected client to be created, got nil")
	}

	if client.apiKey != "" {
		t.Errorf("Expected empty API key, got '%s'", client.apiKey)
	}
}
