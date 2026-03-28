package main

import (
    "os"

    "github.com/Vertinska/weather-app-lab/internal/adapters/weather"
    "github.com/Vertinska/weather-app-lab/internal/pkg/app/cli"
    "github.com/Vertinska/weather-app-lab/internal/pkg/flags"
    "github.com/Vertinska/weather-app-lab/pkg/config"
    "github.com/Vertinska/weather-app-lab/pkg/logger"
)

func main() {
    // Парсим флаги командной строки
    arguments := flags.Parse()
    
    // Создаем логгер с учетом debug режима
    l := logger.New(arguments.Debug)
    
    if arguments.Debug {
        l.Debug("Debug mode enabled")
    }
    
    // Открываем и парсим конфигурационный файл
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
    
    // Получаем провайдера погоды
    wi := getProvider(cfg, l)
    
    // Создаем приложение
    app := cli.New(l, wi, cfg)
    
    // Запускаем
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
        wi = weather.New(l)
    default:
        wi = weather.New(l)
    }
    
    return wi
}
