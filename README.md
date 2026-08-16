# Weather by CEP - Go Cloud Run

A Go application that receives a Brazilian postal code (CEP), identifies the corresponding city, and returns the current weather temperature in Celsius, Fahrenheit, and Kelvin.

## 🔗 Cloud Run URL

**URL**: https://weather-by-cep-179871945135.us-central1.run.app

curl -X POST https://weather-by-cep-179871945135.us-central1.run.app/weather \
  -H "Content-Type: application/json" \
  -d '{"cep":"01310100"}'

---

## 🚀 How to Run Locally

### Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

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

### Response

```json
{
  "temp_c": 25.5,
  "temp_f": 77.9,
  "temp_k": 298.65
}
```