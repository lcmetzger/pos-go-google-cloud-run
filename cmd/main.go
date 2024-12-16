package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func main() {
	godotenv.Load("../.env")

	http.HandleFunc("/", weatherHandler)
	http.ListenAndServe(":3000", nil)
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if len(cep) != 8 {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := getCityByCEP(cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tempC, err := getWeatherByCity(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tempC, tempF, tempK := convertTemperature(tempC)
	response := map[string]float64{
		"temp_C": tempC,
		"temp_F": tempF,
		"temp_K": tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getCityByCEP(cep string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("can not find zipcode")
	}

	var viaCEP ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCEP); err != nil {
		return "", err
	}

	if viaCEP.Localidade == "" {
		return "", fmt.Errorf("can not find zipcode")
	}

	return viaCEP.Localidade, nil
}

func getWeatherByCity(city string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	resp, err := http.Get(fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, city))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("can not find weather data")
	}

	var weatherAPI WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherAPI); err != nil {
		return 0, err
	}

	return weatherAPI.Current.TempC, nil
}

func convertTemperature(tempC float64) (float64, float64, float64) {
	tempF := tempC*1.8 + 32
	tempK := tempC + 273
	return tempC, tempF, tempK
}
