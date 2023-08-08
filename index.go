package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name     string `json:"name"`
	DWeather []struct {
		MainDes string `json:"main"`
	} `json:"weather"`
	Main struct {
		Kelvins float64 `json:"temp"`
		Humid   int     `json:"humidity"`
	} `json:"main"`
}

func getApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Cant read file: %v\n", filename)
		return apiConfigData{}, err
	}

	var c apiConfigData
	err = json.Unmarshal(bytes, &c) // parse json
	if err != nil {
		fmt.Println("Cant get api key:")
		return apiConfigData{}, err
	}
	return c, nil
}

func main() {
	http.HandleFunc("/greeting", greeting)
	http.HandleFunc("/weather/", func(rqw http.ResponseWriter, rq *http.Request) {
		city := strings.SplitN(rq.URL.Path, "/", 3)[2]
		data, err := query(city)
		if err != nil {
			http.Error(rqw, err.Error(), http.StatusInternalServerError)
			return
		}
		rqw.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(rqw).Encode(data)
	})

	http.ListenAndServe(":8080", nil)
}

func query(cityName string) (weatherData, error) {
	apiConfig, err := getApiConfig(".apiConfig")
	if err != nil {
		return weatherData{}, err
	}
	rsp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + cityName)
	if err != nil {
		return weatherData{}, err
	}
	defer rsp.Body.Close()
	var dt weatherData
	if err := json.NewDecoder(rsp.Body).Decode(&dt); err != nil {
		return weatherData{}, err
	}
	return dt, nil

}

func greeting(rqw http.ResponseWriter, rq *http.Request) {
	rqw.Write([]byte("Hi !!"))
}
