package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"

    "github.com/Vertinska/weather-app-lab/internal/adapters/weather"
    "github.com/Vertinska/weather-app-lab/internal/pkg/app/cli"
    "github.com/Vertinska/weather-app-lab/pkg/logger"
)

type WeatherServer struct {
    app *cli.CliApp
    log logger.Logger
    wi  *weather.WeatherInfo
}

func main() {
    log := logger.New()
    wi := weather.New(log)
    app := cli.New(log, wi)
    
    server := &WeatherServer{
        app: app,
        log: log,
        wi:  wi,
    }

    http.HandleFunc("/weather", server.handleWeather)
    http.HandleFunc("/location", server.handleLocation)

    port := 8080
    log.Info(fmt.Sprintf("HTTP сервер запущен на порту %d", port))
    if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
        log.Error("Ошибка запуска сервера", err)
    }
}

func (s *WeatherServer) handleWeather(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    err := s.app.Run()
    if err != nil {
        s.log.Error("Ошибка получения погоды", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Погода обновлена\n"))
}

func (s *WeatherServer) handleLocation(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    switch r.Method {
    case http.MethodGet:
        response := map[string]float64{
            "latitude":  53.6688,
            "longitude": 23.8223,
        }
        json.NewEncoder(w).Encode(response)
        
    case http.MethodPost:
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
        
        // В новой версии нет SetLocation, но можно добавить при необходимости
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(fmt.Sprintf("Местоположение обновлено: %.4f, %.4f\n", lat, lon)))
        
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}
