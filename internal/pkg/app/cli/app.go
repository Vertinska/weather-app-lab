package cli

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    
    "github.com/Vertinska/weather-app-lab/pkg/config"
    "github.com/Vertinska/weather-app-lab/pkg/logger"
)

// CliApp - экспортированная структура (с большой буквы)
type CliApp struct {
    logger      logger.Logger
    locationStg config.LocationStorage
}

// New - конструктор возвращает указатель на CliApp
func New(log logger.Logger, storage config.LocationStorage) *CliApp {
    return &CliApp{
        logger:      log,
        locationStg: storage,
    }
}

func (c *CliApp) Run() error {
    c.logger.Info("Запуск приложения для получения погоды")
    
    // Получаем координаты из хранилища
    latitude, longitude, err := c.locationStg.GetLocation()
    if err != nil {
        c.logger.Errorf("Ошибка при получении координат: %v", err)
        return fmt.Errorf("can't get location: %w", err)
    }
    
    c.logger.Infof("Используем координаты: широта=%.4f, долгота=%.4f", latitude, longitude)

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
    c.logger.Debugf("Отправляем запрос к API: %s", url)

    resp, err := http.Get(url)
    if err != nil {
        c.logger.Errorf("Ошибка при запросе к API: %v", err)
        return fmt.Errorf("can't get weather data from openmeteo: %w", err)
    }
    defer func() {
        if err := resp.Body.Close(); err != nil {
            c.logger.Errorf("Ошибка при закрытии тела ответа: %v", err)
        }
    }()

    c.logger.Debugf("Получен ответ от API, статус: %s", resp.Status)

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        c.logger.Errorf("Ошибка при чтении данных: %v", err)
        return fmt.Errorf("can't read data from response: %w", err)
    }

    c.logger.Debugf("Прочитано %d байт данных", len(data))

    if err := json.Unmarshal(data, &response); err != nil {
        c.logger.Errorf("Ошибка при парсинге JSON: %v", err)
        return fmt.Errorf("can't unmarshal data from response: %w", err)
    }

    message := fmt.Sprintf("Температура воздуха - %.2f градусов цельсия", response.Curr.Temp)
    c.logger.Info(message)
    fmt.Println(message)
    
    return nil
}

// SetLocation устанавливает новое местоположение
func (c *CliApp) SetLocation(lat, lon float64) error {
    c.logger.Infof("Устанавливаем новые координаты: %.4f, %.4f", lat, lon)
    return c.locationStg.SetLocation(lat, lon)
}

// GetLocation возвращает текущее местоположение
func (c *CliApp) GetLocation() (float64, float64, error) {
    return c.locationStg.GetLocation()
}
