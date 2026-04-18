package gui

import (
    "fmt"
    
    "github.com/Vertinska/weather-app-lab/internal/pkg/app/cli"
    "github.com/Vertinska/weather-app-lab/internal/pkg/gui/fyne"
    "github.com/Vertinska/weather-app-lab/pkg/config"
    "github.com/Vertinska/weather-app-lab/pkg/logger"
)

type GUI struct {
    l        *logger.MyLogger
    provider cli.WeatherInfo
    cfg      config.Config
}

func New(l *logger.MyLogger, p cli.WeatherInfo, cfg config.Config) *GUI {
    return &GUI{
        l:        l,
        provider: p,
        cfg:      cfg,
    }
}

func (g *GUI) Run() error {
    g.l.Info("Starting GUI application")
    
    p := fyne.NewP()
    
    window, err := p.CreateWindow("Weather App", fyne.NewWS(400, 300))
    if err != nil {
        return err
    }
    
    textWidget := p.GetTextWidget("Загрузка...")
    window.SetTemperatureWidget(textWidget)
    
    temp, err := g.provider.GetTemperature(g.cfg.L.Lat, g.cfg.L.Long)
    if err != nil {
        g.l.Error("Failed to get temperature", err)
        textWidget.SetText(fmt.Sprintf("Ошибка: %v", err))
    } else {
        textWidget.SetText(fmt.Sprintf("Температура: %.1f°C", temp.Temp))
    }
    
    window.Render()
    p.GetAppRunner().Run()
    
    return nil
}
