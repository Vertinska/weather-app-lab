package pogodaby

import (
    "encoding/json"
    "net/http"
    
    "github.com/Vertinska/weather-app-lab/internal/domain/models"
)

const url = "https://pogoda.by/api/v2/weather-fact?station=26820"

type resp struct {
    Temp float32 `json:"t"`
}

type Logger interface {
    Info(string)
    Debug(string)
    Error(string, error)
    Infof(string, ...interface{})
    Debugf(string, ...interface{})
    Errorf(string, ...interface{})
}

type Pogoda struct {
    l Logger
}

func New(l Logger) *Pogoda {
    return &Pogoda{l: l}
}

func (p *Pogoda) GetTemperature(lat, long float64) (models.TempInfo, error) {
    p.l.Debugf("Fetching weather from pogoda.by for coordinates (%.4f, %.4f)", lat, long)
    
    response, err := http.Get(url)
    if err != nil {
        p.l.Error("can't get data from pogoda.by", err)
        return models.TempInfo{}, err
    }
    defer func() {
        if err := response.Body.Close(); err != nil {
            p.l.Error("can't close response body", err)
        }
    }()

    var r resp
    if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
        p.l.Error("can't decode JSON", err)
        return models.TempInfo{}, err
    }
    
    p.l.Debugf("Successfully got temperature from pogoda.by: %.2f", r.Temp)
    return models.TempInfo{Temp: r.Temp}, nil
}
