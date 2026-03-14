package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "github.com/Vertinska/weather-app-lab/internal/pkg/app/cli"
    "github.com/Vertinska/weather-app-lab/pkg/config"
    "github.com/Vertinska/weather-app-lab/pkg/logger"
)

type WeatherServer struct {
    app *cli.CliApp  // Изменено с App на CliApp
    log logger.Logger
}

func main() {
    log := logger.New(true)
    storage := config.NewFileStorage("")
    app := cli.New(log, storage)  // New возвращает *cliApp
    
    server := &WeatherServer{
        app: app,
        log: log,
    }

    http.HandleFunc("/weather", server.handleWeather)
    http.HandleFunc("/location", server.handleLocation)

    port := 8080
    log.Infof("HTTP сервер запущен на порту %d", port)
    if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
        log.Errorf("Ошибка запуска сервера: %v", err)
    }
}

func (s *WeatherServer) handleWeather(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Получаем погоду через Run
    err := s.app.Run()
    if err != nil {
        s.log.Errorf("Ошибка получения погоды: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Погода обновлена\n"))
}

func (s *WeatherServer) handleLocation(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        // Получить текущее местоположение
        lat, lon, err := s.app.GetLocation()
        if err != nil {
            s.log.Errorf("Ошибка получения местоположения: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        response := map[string]float64{
            "latitude":  lat,
            "longitude": lon,
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
        
    case http.MethodPost:
        // Установить новое местоположение
        latStr := r.URL.Query().Get("lat")
        lonStr := r.URL.Query().Get("lon")
        
        lat, err := strconv.ParseFloat(latStr, 64)
        if err != nil {
            http.Error(w, "Invalid latitude", http.StatusBadRequest)
            return
        }
        
        lon, err := strconv.ParseFloat(lonStr, 64)
        if err != nil {
            http.Error(w, "Invalid longitude", http.StatusBadRequest)
            return
        }
        
        if err := s.app.SetLocation(lat, lon); err != nil {
            s.log.Errorf("Ошибка установки местоположения: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Местоположение обновлено\n"))
        
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
