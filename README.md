# Weather by CEP - Go Cloud Run

Go HTTP service that takes a Brazilian postal code (CEP), resolves it to a city via [ViaCEP](https://viacep.com.br/), and returns the current temperature in Celsius, Fahrenheit and Kelvin from [WeatherAPI](https://www.weatherapi.com/). Containerized with Docker and deployed on Google Cloud Run.

## 🔗 Cloud Run URL

**URL**: https://weather-by-cep-261668892219.us-central1.run.app

```bash
curl -X POST https://weather-by-cep-261668892219.us-central1.run.app/weather \
  -H "Content-Type: application/json" \
  -d '{"cep":"01310100"}'
```

---

## 🚀 How to Run Locally

### Prerequisites

- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
- [Go 1.26+](https://go.dev/dl/) — only to run `make test`
- A free [WeatherAPI](https://www.weatherapi.com/) key
- [jq](https://jqlang.github.io/jq/) — optional, used to format the `make test-*` output

### Setup

```bash
cp .env.sample .env   # then set WEATHER_API_KEY
```

### Running the Application with Docker

```bash
# 1. Build Docker images
make build

# 2. Start containers
make up

# 3. Application will be available at
# http://localhost:8080
```

### Running Tests

```bash
# Run all unit tests
make test

# Test with valid CEP (São Paulo - SP)
make test-valid

# Test with invalid CEP
make test-invalid

# Test with non-existent CEP
make test-notfound
```

### Stopping the Containers

```bash
make down
```

---

## API

### Endpoint

`POST /weather`

### Request

```json
{
  "cep": "01310100"
}
```

### Response — `200 OK`

```json
{
  "temp_C": 25.5,
  "temp_F": 77.9,
  "temp_K": 298.5
}
```

### Errors

| Status | Condition | Body |
|---|---|---|
| `422` | CEP is not 8 digits, or has invalid characters | `{"message":"invalid zipcode"}` |
| `404` | CEP is well-formed but does not exist | `{"message":"can not find zipcode"}` |

Conversions: `F = C × 1.8 + 32` · `K = C + 273`
