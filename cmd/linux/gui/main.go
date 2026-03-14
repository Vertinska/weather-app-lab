package main

import (
    "fmt"
    "time"

    "github.com/Vertinska/weather-app-lab/internal/pkg/app/cli"
    "github.com/Vertinska/weather-app-lab/pkg/config"
    "github.com/Vertinska/weather-app-lab/pkg/logger"
    "github.com/getlantern/systray"
)

func main() {
    log := logger.New(false)
    storage := config.NewFileStorage("")
    app := cli.New(log, storage)
    
    // Запускаем системный трей
    systray.Run(func() {
        onReady(app)
    }, func() {
        onExit()
    })
}

func onReady(app *cli.CliApp) {
    systray.SetTitle("Weather App")
    systray.SetTooltip("Погода")
    
    // Создаем меню
    mWeather := systray.AddMenuItem("Обновить погоду", "Показать текущую погоду")
    systray.AddSeparator()
    
    // Предустановленные города
    mLocationGrodno := systray.AddMenuItem("Гродно", "Установить местоположение - Гродно")
    mLocationMinsk := systray.AddMenuItem("Минск", "Установить местоположение - Минск")
    mLocationBrest := systray.AddMenuItem("Брест", "Установить местоположение - Брест")
    mLocationVitebsk := systray.AddMenuItem("Витебск", "Установить местоположение - Витебск")
    mLocationGomel := systray.AddMenuItem("Гомель", "Установить местоположение - Гомель")
    mLocationMogilev := systray.AddMenuItem("Могилев", "Установить местоположение - Могилев")
    
    systray.AddSeparator()
    mQuit := systray.AddMenuItem("Выход", "Закрыть приложение")
    
    // Таймер для автоматического обновления каждый час
    ticker := time.NewTicker(1 * time.Hour)
    go func() {
        for {
            updateWeather(app)
            <-ticker.C
        }
    }()
    
    // Обработчики меню
    go func() {
        for {
            select {
            case <-mWeather.ClickedCh:
                updateWeather(app)
            case <-mLocationGrodno.ClickedCh:
                app.SetLocation(53.6688, 23.8223)
                updateWeather(app)
            case <-mLocationMinsk.ClickedCh:
                app.SetLocation(53.9045, 27.5615)
                updateWeather(app)
            case <-mLocationBrest.ClickedCh:
                app.SetLocation(52.0976, 23.7341)
                updateWeather(app)
            case <-mLocationVitebsk.ClickedCh:
                app.SetLocation(55.1904, 30.2049)
                updateWeather(app)
            case <-mLocationGomel.ClickedCh:
                app.SetLocation(52.4345, 30.9754)
                updateWeather(app)
            case <-mLocationMogilev.ClickedCh:
                app.SetLocation(53.9168, 30.3449)
                updateWeather(app)
            case <-mQuit.ClickedCh:
                systray.Quit()
                return
            }
        }
    }()
    
    // Показываем погоду при запуске
    go func() {
        time.Sleep(1 * time.Second)
        updateWeather(app)
    }()
}

func updateWeather(app *cli.CliApp) {
    if err := app.Run(); err != nil {
        fmt.Println("Ошибка:", err)
        systray.SetTooltip("❌ Ошибка получения погоды")
    } else {
        systray.SetTooltip("✅ Погода обновлена")
    }
}

func onExit() {
    fmt.Println("Приложение завершено")
}
