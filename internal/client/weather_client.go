package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WeatherClient struct {
	WeatherAPI    string
	WeatherClient *http.Client
}

// WeatherResponse представляет полный ответ API погоды
type WeatherResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
	Forecast Forecast `json:"forecast"`
}

// Location содержит информацию о местоположении
type Location struct {
	Name      string  `json:"name"`
	Region    string  `json:"region"`
	Country   string  `json:"country"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Timezone  string  `json:"tz_id"`
	Localtime string  `json:"localtime"`
}

// Current содержит текущие погодные данные
type Current struct {
	LastUpdated string    `json:"last_updated"`
	TempC       float64   `json:"temp_c"`
	IsDay       int       `json:"is_day"`
	Condition   Condition `json:"condition"`
	WindKPH     float64   `json:"wind_kph"`
	WindDir     string    `json:"wind_dir"`
}

// Forecast содержит прогноз погоды
type Forecast struct {
	ForecastDays []ForecastDay `json:"forecastday"`
}

// ForecastDay представляет прогноз на один день
type ForecastDay struct {
	Date string `json:"date"`
	Day  Day    `json:"day"`
}

// Day содержит дневные погодные данные
type Day struct {
	MaxTempC     float64   `json:"maxtemp_c"`
	MinTempC     float64   `json:"mintemp_c"`
	AvgTempC     float64   `json:"avgtemp_c"`
	MaxWindKPH   float64   `json:"maxwind_kph"`
	TotalSnowCM  float64   `json:"totalsnow_cm"`
	WillItRain   int       `json:"daily_will_it_rain"`
	ChanceOfRain int       `json:"daily_chance_of_rain"`
	WillItSnow   int       `json:"daily_will_it_snow"`
	ChanceOfSnow int       `json:"daily_chance_of_snow"`
	Condition    Condition `json:"condition"`
	UV           float64   `json:"uv"`
}

// Condition описывает погодные условия
type Condition struct {
	Text string `json:"text"`
}

func NewWeatherClient(apikey string) *WeatherClient {
	return &WeatherClient{
		WeatherAPI:    apikey,
		WeatherClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *WeatherClient) GetCurrentWeather(location string, day int) (string, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=14&aqi=no&alerts=no",
		c.WeatherAPI, location)

	resp, err := c.WeatherClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("ошибка при запросе к API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API вернул ошибку: %v", resp.StatusCode)
	}

	var weatherResponse WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return "", fmt.Errorf("ошибка при чтение json: %v", err)
	}

	result := fmt.Sprintf(
		//	"Погода в %s (%s): %.1f°C, %s, Ветер: %.1f км/ч, Влажность: %d%%",
		"Погода в %s (%s) на %s: \nМаксимальная температура: %.1f°C.\nМинимальная температура: %.1f°C\nСредняя температура :%1.f°C\nСкорость ветра: %.1fкм/ч.\nВероятность дождя (в процентах): %v \nОписание: %s",
		weatherResponse.Location.Name,
		weatherResponse.Location.Country,
		weatherResponse.Forecast.ForecastDays[day].Date,
		weatherResponse.Forecast.ForecastDays[day].Day.MaxTempC,
		weatherResponse.Forecast.ForecastDays[day].Day.MinTempC,
		weatherResponse.Forecast.ForecastDays[day].Day.AvgTempC,
		weatherResponse.Forecast.ForecastDays[day].Day.MaxWindKPH,
		weatherResponse.Forecast.ForecastDays[day].Day.ChanceOfRain,
		weatherResponse.Forecast.ForecastDays[day].Day.Condition.Text,
	)

	return result, nil
}
