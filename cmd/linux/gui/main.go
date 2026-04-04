package main

import (
    "fmt"
    "time"

    pogodaby "github.com/Vertinska/weather-app-lab/internal/adapters/pogoda_by"
    "github.com/Vertinska/weather-app-lab/internal/adapters/weather"
    "github.com/Vertinska/weather-app-lab/internal/pkg/app/cli"
    "github.com/Vertinska/weather-app-lab/pkg/config"
    "github.com/Vertinska/weather-app-lab/pkg/logger"
    "github.com/getlantern/systray"
)

func main() {
    l := logger.New(false)
    
    // Загружаем конфиг
    cfg, err := config.LoadConfigFromFile("./config/config.yaml")
    if err != nil {
        l.Errorf("Failed to load config: %v", err)
        return
    }
    
    // Выбираем провайдера
    var wi cli.WeatherInfo
    if cfg.P.Type == "pogoda" {
        wi = pogodaby.New(l)
    } else {
        wi = weather.NewWithCacheTTL(l, 600*1000000000)
    }
    
    app := cli.New(l, wi, cfg)
    
    // Запускаем системный трей
    systray.Run(func() {
        onReady(app, l)
    }, func() {
        onExit()
    })
}

func onReady(app *cli.CliApp, l logger.Logger) {
    systray.SetTitle("Weather App")
    systray.SetTooltip("Погода")
    
    // Создаем меню
    mWeather := systray.AddMenuItem("Обновить погоду", "Показать текущую погоду")
    systray.AddSeparator()
    mQuit := systray.AddMenuItem("Выход", "Закрыть приложение")
    
    // Обработчики меню
    go func() {
        for {
            select {
            case <-mWeather.ClickedCh:
                if err := app.Run(); err != nil {
                    l.Error("Ошибка получения погоды", err)
                    systray.SetTooltip("❌ Ошибка")
                } else {
                    systray.SetTooltip("✅ Погода обновлена")
                }
            case <-mQuit.ClickedCh:
                systray.Quit()
                return
            }
        }
    }()
    
    // Показываем погоду при запуске
    go func() {
        time.Sleep(1 * time.Second)
        app.Run()
    }()
}

func onExit() {
    fmt.Println("Приложение завершено")
}
