package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"

    "github.com/Vertinska/weather-app-lab/internal/pkg/app/cli"
    "github.com/Vertinska/weather-app-lab/pkg/config"
    "github.com/Vertinska/weather-app-lab/pkg/logger"
)

func main() {
    debugMode := flag.Bool("debug", false, "Включить debug режим")
    setLocation := flag.String("set-location", "", "Установить местоположение (формат: широта,долгота)")
    configPath := flag.String("config", "", "Путь к конфигурационному файлу")
    flag.Parse()

    l := logger.New()
    
    if *debugMode {
        l.Debug("Debug режим включен")
    }

    cfg, err := loadConfig(*configPath)
    if err != nil {
        l.Error("Ошибка при загрузке конфигурации", err)
        os.Exit(1)
    }

    storage, err := createLocationStorage(cfg, l)
    if err != nil {
        l.Error("Ошибка при создании хранилища", err)
        os.Exit(1)
    }

    app := cli.New(l, storage)

    if *setLocation != "" {
        var lat, lon float64
        n, err := fmt.Sscanf(*setLocation, "%f,%f", &lat, &lon)
        if err != nil || n != 2 {
            l.Error("Неверный формат местоположения", fmt.Errorf("используйте: широта,долгота"))
            os.Exit(1)
        }
        
        if err := app.SetLocation(lat, lon); err != nil {
            l.Error("Ошибка при установке местоположения", err)
            os.Exit(1)
        }
        l.Info("Местоположение успешно сохранено")
        return
    }

    if err := app.Run(); err != nil {
        l.Error("Ошибка при выполнении приложения", err)
        os.Exit(1)
    }
    
    os.Exit(0)
}

func loadConfig(path string) (*config.Config, error) {
    if path == "" {
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return nil, err
        }
        path = filepath.Join(homeDir, ".weather-app", "config.json")
    }
    return config.LoadConfig(path)
}

func createLocationStorage(cfg *config.Config, l logger.Logger) (config.LocationStorage, error) {
    switch cfg.StorageType {
    case "file":
        l.Debug("Используем файловое хранилище")
        return config.NewFileStorage(cfg.FilePath), nil
    case "memory":
        l.Debug("Используем хранилище в памяти")
        return config.NewMemoryStorage(), nil
    default:
        l.Debug(fmt.Sprintf("Неизвестный тип хранилища '%s', используем file", cfg.StorageType))
        return config.NewFileStorage(""), nil
    }
}
