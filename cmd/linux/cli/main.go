package main

import (
    "os"
    "time"

    pogodaby "github.com/Vertinska/weather-app-lab/internal/adapters/pogoda_by"
    "github.com/Vertinska/weather-app-lab/internal/adapters/weather"
    "github.com/Vertinska/weather-app-lab/internal/pkg/app/cli"
    "github.com/Vertinska/weather-app-lab/internal/pkg/flags"
    "github.com/Vertinska/weather-app-lab/pkg/config"
    "github.com/Vertinska/weather-app-lab/pkg/logger"
)

func main() {
    arguments := flags.Parse()
    
    l := logger.New(arguments.Debug)
    
    if arguments.Debug {
        l.Debug("Debug mode enabled")
    }
    
    r, err := os.Open(arguments.Path)
    if err != nil {
        l.Errorf("Failed to open config file: %v", err)
        os.Exit(1)
    }
    defer r.Close()
    
    cfg, err := config.Parse(r)
    if err != nil {
        l.Errorf("Failed to parse config: %v", err)
        os.Exit(1)
    }
    
    l.Debugf("Config loaded successfully")
    l.Debugf("Using coordinates: %.4f, %.4f", cfg.L.Lat, cfg.L.Long)
    l.Debugf("Provider type: %s", cfg.P.Type)
    l.Debugf("Cache TTL: %d seconds", cfg.P.CacheTTL)
    
    wi := getProvider(cfg, l)
    app := cli.New(l, wi, cfg)
    
    err = app.Run()
    if err != nil {
        l.Error("Some error", err)
        os.Exit(1)
    }
    
    os.Exit(0)
}

func getProvider(cfg config.Config, l cli.Logger) cli.WeatherInfo {
    var wi cli.WeatherInfo
    
    switch cfg.P.Type {
    case "open-meteo":
        ttl := time.Duration(cfg.P.CacheTTL) * time.Second
        wi = weather.NewWithCacheTTL(l, ttl)
    case "pogoda":
        wi = pogodaby.New(l)
    default:
        wi = weather.NewWithCacheTTL(l, 10*time.Minute)
    }
    
    return wi
}
