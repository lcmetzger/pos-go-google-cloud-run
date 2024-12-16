package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
)

func TestWeatherHandler_Success(t *testing.T) {
	godotenv.Load("../.env")

	req, err := http.NewRequest("GET", "/?cep=89040115", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(weatherHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

}

func TestWeatherHandler_InvalidCEP(t *testing.T) {
	req, err := http.NewRequest("GET", "/?cep=123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(weatherHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}

	expected := "invalid zipcode\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestWeatherHandler_CEPNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/?cep=00000000", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(weatherHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	expected := "can not find zipcode\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
