package providers

import (
    "time"

    pogodaby "github.com/Vertinska/weather-app-lab/internal/adapters/pogoda_by"
    "github.com/Vertinska/weather-app-lab/internal/adapters/weather"
    "github.com/Vertinska/weather-app-lab/internal/pkg/app/cli"
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

func GetProvider(cfg config.Config, l Logger) cli.WeatherInfo {
    switch cfg.P.Type {
    case "open-meteo":
        ttl := time.Duration(cfg.P.CacheTTL) * time.Second
        return weather.NewWithCacheTTL(l, ttl)
    case "pogoda":
        return pogodaby.New(l)
    default:
        ttl := time.Duration(cfg.P.CacheTTL) * time.Second
        return weather.NewWithCacheTTL(l, ttl)
    }
}
