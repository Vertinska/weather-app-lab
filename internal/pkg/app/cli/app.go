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
    GetTemperature(float64, float64) (models.TempInfo, error)
}

type CliApp struct {
    l    Logger
    wi   WeatherInfo
    conf config.Config
}

func New(l Logger, wi WeatherInfo, c config.Config) *CliApp {
    return &CliApp{
        l:    l,
        wi:   wi,
        conf: c,
    }
}

func (c *CliApp) Run() error {
    tempInfo, err := c.wi.GetTemperature(c.conf.L.Lat, c.conf.L.Long)
    if err != nil {
        c.l.Error("can't get temp info", err)
        return err
    }
    fmt.Printf(
        "Температура воздуха - %.2f градусов цельсия\n",
        tempInfo.Temp,
    )
    return nil
}

func (c *CliApp) SetLocation(lat, lon float64) error {
    c.l.Debugf("SetLocation called with: %.4f, %.4f", lat, lon)
    return nil
}

func (c *CliApp) GetLocation() (float64, float64, error) {
    return c.conf.L.Lat, c.conf.L.Long, nil
}
