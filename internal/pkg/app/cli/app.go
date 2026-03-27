package cli

import (
    "fmt"
    
    "github.com/Vertinska/weather-app-lab/internal/domain/models"
)

type Logger interface {
    Info(string)
    Debug(string)
    Error(string, error)
}

type WeatherInfo interface {
    GetTemperature(float64, float64) models.TempInfo
}

type CliApp struct {
    l  Logger
    wi WeatherInfo
}

func New(l Logger, wi WeatherInfo) *CliApp {
    return &CliApp{
        l:  l,
        wi: wi,
    }
}

func (c *CliApp) Run() error {
    // Координаты Гродно
    latitude := 53.6688
    longitude := 23.8223
    
    temp := c.wi.GetTemperature(latitude, longitude)
    
    fmt.Printf(
        "Температура воздуха - %.2f градусов цельсия\n",
        temp.Temp,
    )
    return nil
}

// Дополнительные методы для совместимости с предыдущими версиями
func (c *CliApp) SetLocation(lat, lon float64) error {
    c.l.Debug(fmt.Sprintf("SetLocation called with: %.4f, %.4f", lat, lon))
    return nil
}

func (c *CliApp) GetLocation() (float64, float64, error) {
    return 53.6688, 23.8223, nil
}
