package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var apiKey string

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	return os.Getenv(key)
}

func main() {
	startNow := time.Now()

	apiKey = getEnvVariable("API_KEY")

	cities := []string{"London", "Glasgow", "Liverpool", "Perth", "New+York", "York", "Paris", "Tokyo", "Kyoto", "Seoul"}

	ch := make(chan string)
	var wg sync.WaitGroup

	for _, city := range cities {
		wg.Add(1)
		go fetchWeather(city, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		fmt.Println(result)
	}

	fmt.Println("This operation took:", time.Since(startNow))
}

func fetchWeather(city string, ch chan<- string, wg *sync.WaitGroup) interface{} {
	var WeatherData struct {
		Data []struct {
			Temp     float64 `json:"temp"`
			CityName string  `json:"city_name"`
		} `json:"data"`
	}

	defer wg.Done()

	url := fmt.Sprintf("https://api.weatherbit.io/v2.0/current?city=%s&key=%s", city, apiKey)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("Error fetching weather for %s: %s\n", city, err)
		return WeatherData
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&WeatherData); err != nil {
		fmt.Printf("Error fetching weather data for %s: %s\n", city, err)
		return WeatherData
	}

	ch <- fmt.Sprintf("This is the data - city: %s; temp: %f;", WeatherData.Data[0].CityName, WeatherData.Data[0].Temp)

	return WeatherData
}
