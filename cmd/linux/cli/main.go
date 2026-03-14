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
    // Парсим флаги командной строки
    debugMode := flag.Bool("debug", false, "Включить debug режим")
    setLocation := flag.String("set-location", "", "Установить местоположение (формат: широта,долгота)")
    configPath := flag.String("config", "", "Путь к конфигурационному файлу")
    flag.Parse()

    // Создаем логгер
    log := logger.New(*debugMode)
    
    log.Info("Старт приложения")
    if *debugMode {
        log.Debug("Debug режим включен")
    }

    // Загружаем конфигурацию
    cfg, err := loadConfig(*configPath)
    if err != nil {
        log.Error("Ошибка при загрузке конфигурации: " + err.Error())
        os.Exit(1)
    }
    log.Debugf("Загружена конфигурация: тип хранилища = %s", cfg.StorageType)

    // Создаем хранилище местоположения
    storage, err := createLocationStorage(cfg, log)
    if err != nil {
        log.Error("Ошибка при создании хранилища: " + err.Error())
        os.Exit(1)
    }

    // Создаем приложение
    app := cli.New(log, storage)

    // Если указана команда установки местоположения
    if *setLocation != "" {
        var lat, lon float64
        n, err := fmt.Sscanf(*setLocation, "%f,%f", &lat, &lon)
        if err != nil || n != 2 {
            log.Error("Неверный формат местоположения. Используйте: широта,долгота")
            os.Exit(1)
        }
        
        if err := app.SetLocation(lat, lon); err != nil {
            log.Error("Ошибка при установке местоположения: " + err.Error())
            os.Exit(1)
        }
        log.Info("Местоположение успешно сохранено")
        return
    }

    // Запускаем приложение
    if err := app.Run(); err != nil {
        log.Error("Ошибка при выполнении приложения: " + err.Error())
        os.Exit(1)
    }
    
    log.Info("Приложение успешно завершило работу")
}

func loadConfig(path string) (*config.Config, error) {
    if path == "" {
        // Используем стандартный путь
        homeDir, err := os.UserHomeDir()
        if err != nil {
            return nil, err
        }
        path = filepath.Join(homeDir, ".weather-app", "config.json")
    }
    
    return config.LoadConfig(path)
}

func createLocationStorage(cfg *config.Config, log logger.Logger) (config.LocationStorage, error) {
    switch cfg.StorageType {
    case "file":
        log.Debug("Используем файловое хранилище")
        return config.NewFileStorage(cfg.FilePath), nil
    case "memory":
        log.Debug("Используем хранилище в памяти")
        return config.NewMemoryStorage(), nil
    default:
        log.Debugf("Неизвестный тип хранилища '%s', используем file", cfg.StorageType)
        return config.NewFileStorage(""), nil
    }
}
