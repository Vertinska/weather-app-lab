package main

import (
    "os"

    "github.com/Vertinska/weather-app-lab/internal/adapters/weather"
    "github.com/Vertinska/weather-app-lab/internal/pkg/app/cli"
    "github.com/Vertinska/weather-app-lab/pkg/logger"
)

func main() {
    l := logger.New()
    wi := weather.New(l)
    app := cli.New(l, wi)
    
    err := app.Run()
    if err != nil {
        l.Error("Some error", err)
        os.Exit(1)
    }
    
    os.Exit(0)
}
