package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"regexp"

	"fullcyclelabs.com/weather-app/internal/clients/viacep"
	"fullcyclelabs.com/weather-app/internal/clients/weatherapi"
)

const (
	ContentType     string = "Content-Type"
	ApplicationJson string = "application/json"
)

var ErrZipCodeNotFound = errors.New("zipcode not found")

type WeatherRequest struct {
	CEP string `json:"cep"`
}

type WeatherResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type WeatherHandler struct {
	viacepClient     *viacep.Client
	weatherapiClient *weatherapi.Client
}

func NewWeatherHandler(
	viacepClient *viacep.Client,
	weatherapiClient *weatherapi.Client,
) *WeatherHandler {
	return &WeatherHandler{
		viacepClient:     viacepClient,
		weatherapiClient: weatherapiClient,
	}
}

func (h *WeatherHandler) HandleWeatherRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		w.Header().Set(ContentType, ApplicationJson)
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "method not allowed"})
		return
	}

	var req WeatherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set(ContentType, ApplicationJson)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid request format"})
		return
	}

	if !isValidCEP(req.CEP) {
		w.Header().Set(ContentType, ApplicationJson)
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid zipcode"})
		return
	}

	city, err := h.getCityByCEP(ctx, req.CEP)
	if err != nil {
		w.Header().Set(ContentType, ApplicationJson)
		if errors.Is(err, ErrZipCodeNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: "can not find zipcode"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		}
		return
	}

	tempC, err := h.getTemperatureByCity(ctx, city)
	if err != nil {
		w.Header().Set(ContentType, ApplicationJson)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		return
	}

	tempF := celsiusToFahrenheit(tempC)
	tempK := celsiusToKelvin(tempC)

	response := WeatherResponse{
		TempC: math.Round(tempC*10) / 10,
		TempF: math.Round(tempF*10) / 10,
		TempK: math.Round(tempK*10) / 10,
	}

	w.Header().Set(ContentType, ApplicationJson)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func isValidCEP(cep string) bool {
	matched, _ := regexp.MatchString(`^\d{8}$`, cep)
	return matched
}

func (h *WeatherHandler) getCityByCEP(ctx context.Context, cep string) (string, error) {
	client := h.viacepClient

	res, err := client.FetchByCep(ctx, cep)
	if err != nil {
		return "", fmt.Errorf("failed to call ViaCEP: %w", err)
	}

	if res.Error != "" {
		return "", ErrZipCodeNotFound
	}

	if res.City == "" {
		return "", ErrZipCodeNotFound
	}

	return res.City, nil
}

func (h *WeatherHandler) getTemperatureByCity(ctx context.Context, city string) (float64, error) {
	client := h.weatherapiClient

	res, err := client.FetchByCity(ctx, city)
	if err != nil {
		return 0, fmt.Errorf("failed to call WeatherAPI: %w", err)
	}

	return res.Current.TempC, nil
}

func celsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

func celsiusToKelvin(c float64) float64 {
	return c + 273
}
