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

type WeatherResponse struct {
	Location struct {
		Name      string `json:"name"`
		Region    string `json:"region"`
		Country   string `json:"country"`
		LocalTime string `json:"localtime"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		WindKph float64 `json:"wind_kph"`
		WindDir string  `json:"wind_dir"`
	} `json:"current"`
	Forecast struct {
		Data string `json:"date"`
		Day  struct {
			MaxTempC          float64 `json:"maxtemp_c"`
			MinTempC          float64 `json:"mintemp_c"`
			AvgTempC          float64 `json:"avgtemp_c"`
			MaxWindKph        float64 `json:"maxwind_kph"`
			TotalSnowCM       float64 `json:"totalsnow_cm"`
			DailyWillItRain   float64 `json:"daily_will_it_rain"`
			DailyChanceOfRain float64 `json:"daily_chance_of_rain"`
			Condition         struct {
				Text string `json:"text"`
			} `json:"condition"`
			Uv float64 `json:"uv"`
		} `json:"day"`
	} `json:"forecast"`
}

func NewWeatherClient(apikey string) *WeatherClient {
	return &WeatherClient{
		WeatherAPI:    apikey,
		WeatherClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *WeatherClient) GetCurrentWeather(location string) (string, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no",
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
		"Погода в %s (%s) на %s: \nТемпература: %.1f°C.",
		weatherResponse.Location.Name,
		weatherResponse.Location.Country,
		weatherResponse.Location.LocalTime[:10],
		weatherResponse.Current.TempC,
	)

	return result, nil
}
