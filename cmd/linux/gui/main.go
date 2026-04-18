package main

import (
    "os"

    "github.com/Vertinska/weather-app-lab/internal/pkg/app/gui"
    "github.com/Vertinska/weather-app-lab/pkg/config"
    "github.com/Vertinska/weather-app-lab/internal/pkg/flags"
    "github.com/Vertinska/weather-app-lab/internal/pkg/providers"
    "github.com/Vertinska/weather-app-lab/pkg/logger"
)

func main() {
    arguments := flags.Parse()

    r, err := os.Open(arguments.Path)
    if err != nil {
        panic(err)
    }
    defer r.Close()

    cfg, err := config.Parse(r)
    if err != nil {
        panic(err)
    }

    l := logger.New(arguments.Debug)
    provider := providers.GetProvider(cfg, l)
    
    g := gui.New(l, provider, cfg)
    
    err = g.Run()
    if err != nil {
        panic(err)
    }
}
