package cli

import (
    "fmt"
    
    "github.com/Vertinska/weather-app-lab/internal/domain/models"
    "github.com/Vertinska/weather-app-lab/pkg/config"
)

type Logger interface {
    Info(string)
    Debug(string)
    Error(string, error)
    Infof(string, ...interface{})
    Debugf(string, ...interface{})
    Errorf(string, ...interface{})
}

type WeatherInfo interface {
    GetTemperature(float64, float64) models.TempInfo
}

type CliApp struct {
    l      Logger
    wi     WeatherInfo
    cfg    config.Config
}

func New(l Logger, wi WeatherInfo, cfg config.Config) *CliApp {
    return &CliApp{
        l:   l,
        wi:  wi,
        cfg: cfg,
    }
}

func (c *CliApp) Run() error {
    // Используем координаты из конфигурации
    latitude := c.cfg.L.Lat
    longitude := c.cfg.L.Long
    
    temp := c.wi.GetTemperature(latitude, longitude)
    
    fmt.Printf(
        "Температура воздуха - %.2f градусов цельсия\n",
        temp.Temp,
    )
    return nil
}

// Дополнительные методы для совместимости
func (c *CliApp) SetLocation(lat, lon float64) error {
    c.l.Debugf("SetLocation called with: %.4f, %.4f", lat, lon)
    return nil
}

func (c *CliApp) GetLocation() (float64, float64, error) {
    return c.cfg.L.Lat, c.cfg.L.Long, nil
}
