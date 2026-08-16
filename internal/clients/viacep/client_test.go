package viacep

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViaCEPResponse_Parsing(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expectCity  string
		expectError string
	}{
		{
			name:       "Valid response with city",
			jsonData:   `{"localidade":"São Paulo","uf":"SP","logradouro":"Avenida Paulista"}`,
			expectCity: "São Paulo",
		},
		{
			name:        "Response with error flag",
			jsonData:    `{"erro":"true"}`,
			expectError: "true",
		},
		{
			name:       "Valid response",
			jsonData:   `{"localidade":"Rio de Janeiro","erro":""}`,
			expectCity: "Rio de Janeiro",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp ViaCEPResponse
			if err := json.Unmarshal([]byte(tt.jsonData), &resp); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if tt.expectCity != "" && resp.City != tt.expectCity {
				t.Errorf("Expected city '%s', got '%s'", tt.expectCity, resp.City)
			}

			if tt.expectError != "" && resp.Error != tt.expectError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectError, resp.Error)
			}
		})
	}
}

func TestViaCEPClient_Structure(t *testing.T) {
	client := NewClient(http.DefaultClient)

	if client == nil {
		t.Fatalf("Expected client to be created, got nil")
	}

	if client.client == nil {
		t.Errorf("Expected client.client to be initialized")
	}
}

func TestViaCEPClient_FetchByCep_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"localidade":"São Paulo","uf":"SP"}`))
	}))
	defer server.Close()

	httpClient := server.Client()
	_ = NewClient(httpClient)

	// Note: In a real test, you might mock the HTTP client completely
	// or modify the client to accept a base URL for testing
}

func TestViaCEPResponse_FieldMapping(t *testing.T) {
	jsonStr := `{"localidade":"São Paulo","erro":""}`

	var resp ViaCEPResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if resp.City != "São Paulo" {
		t.Errorf("Expected City field to be 'São Paulo', got '%s'", resp.City)
	}
}

func TestViaCEPClient_Integration(t *testing.T) {
	// This is a basic integration test that can be run against the real API
	// Comment this out if you don't want to make real API calls during testing
	t.Skip("Skipping integration test - requires external API access")

	client := NewClient(http.DefaultClient)
	ctx := context.Background()

	resp, err := client.FetchByCep(ctx, "01310100")
	if err != nil {
		t.Errorf("FetchByCep returned error: %v", err)
		return
	}

	if resp == nil {
		t.Errorf("Expected response, got nil")
		return
	}

	if resp.City == "" {
		t.Errorf("Expected city to be populated")
	}
}

func TestViaCEPClient_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// The actual test would depend on how the client handles HTTP errors
	// This is a template for how you might structure such tests
}
