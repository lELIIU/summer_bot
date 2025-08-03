package service

import (
	"summer_bot/internal/client"
)

type WeatherService struct {
	Client *client.WeatherClient
}

func NewWeatherService(apiKey string) *WeatherService {
	return &WeatherService{
		Client: client.NewWeatherClient(apiKey),
	}
}

func (s *WeatherService) GetWeather(city string) (string, error) {
	return s.Client.GetCurrentWeather(city)
}
