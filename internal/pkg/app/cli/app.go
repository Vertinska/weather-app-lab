package cli

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    
    "github.com/Vertinska/weather-app-lab/pkg/config"
    "github.com/Vertinska/weather-app-lab/pkg/logger"
)

type CliApp struct {
    l           logger.Logger
    locationStg config.LocationStorage
}

func New(l logger.Logger, storage config.LocationStorage) *CliApp {
    return &CliApp{
        l:           l,
        locationStg: storage,
    }
}

func (c *CliApp) Run() error {
    c.l.Info("Запуск приложения для получения погоды")
    
    // Получаем координаты из хранилища
    latitude, longitude, err := c.locationStg.GetLocation()
    if err != nil {
        c.l.Error("Ошибка при получении координат", err)
        return fmt.Errorf("can't get location: %w", err)
    }
    
    c.l.Debug(fmt.Sprintf("Используем координаты: широта=%.4f, долгота=%.4f", latitude, longitude))

    type Current struct {
        Temp float32 `json:"temperature_2m"`
    }

    type Response struct {
        Curr Current `json:"current"`
    }

    var response Response

    params := fmt.Sprintf(
        "latitude=%f&longitude=%f&current=temperature_2m",
        latitude,
        longitude,
    )

    url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?%s", params)
    c.l.Debug(fmt.Sprintf("url was generated success - %s", url))

    resp, err := http.Get(url)
    if err != nil {
        c.l.Error("can't get weather data", err)
        customErr := errors.New("can't get weather data from openmeteo")
        return errors.Join(customErr, err)
    }
    defer func() {
        if err := resp.Body.Close(); err != nil {
            c.l.Error("can't close body", err)
        }
    }()

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        c.l.Error("can't read data from body", err)
        customErr := errors.New("can't read data from response")
        return errors.Join(customErr, err)
    }

    c.l.Debug(fmt.Sprintf("data was readed successfully size - %d", len(data)))

    if err := json.Unmarshal(data, &response); err != nil {
        c.l.Error("can't unmarshal json data", err)
        customErr := errors.New("can't unmarshal data from response")
        return errors.Join(customErr, err)
    }

    message := fmt.Sprintf("Температура воздуха - %.2f градусов цельсия", response.Curr.Temp)
    c.l.Info(message)
    fmt.Println(message)
    
    return nil
}

func (c *CliApp) SetLocation(lat, lon float64) error {
    c.l.Debug(fmt.Sprintf("Устанавливаем новые координаты: %.4f, %.4f", lat, lon))
    return c.locationStg.SetLocation(lat, lon)
}

func (c *CliApp) GetLocation() (float64, float64, error) {
    return c.locationStg.GetLocation()
}
