package weather

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "time"
    
    "github.com/Vertinska/weather-app-lab/internal/domain/models"
    "github.com/Vertinska/weather-app-lab/internal/pkg/cache"
)

const apiURL = "https://api.open-meteo.com/v1/forecast"

type Logger interface {
    Info(string)
    Debug(string)
    Error(string, error)
    Infof(string, ...interface{})
    Debugf(string, ...interface{})
    Errorf(string, ...interface{})
}

type current struct {
    Temp float32 `json:"temperature_2m"`
}

type response struct {
    Curr current `json:"current"`
}

type WeatherInfo struct {
    c        current
    l        Logger
    isLoaded bool
    cache    *cache.FileCache
    cacheTTL time.Duration
}

func New(l Logger) *WeatherInfo {
    return NewWithCacheTTL(l, 10*time.Minute)
}

func NewWithCacheTTL(l Logger, ttl time.Duration) *WeatherInfo {
    homeDir, _ := os.UserHomeDir()
    cacheDir := filepath.Join(homeDir, ".weather-app-cache")
    
    return &WeatherInfo{
        l:        l,
        cache:    cache.NewFileCache(cacheDir, ttl),
        cacheTTL: ttl,
    }
}

func (wi *WeatherInfo) getCacheKey(lat, long float64) string {
    return fmt.Sprintf("weather_%.6f_%.6f", lat, long)
}

func (wi *WeatherInfo) getWeatherInfo(lat, long float64) error {
    cacheKey := wi.getCacheKey(lat, long)
    
    if cached, found := wi.cache.Get(cacheKey); found {
        if temp, ok := cached.(float64); ok {
            wi.c = current{Temp: float32(temp)}
            wi.isLoaded = true
            wi.l.Debugf("Data loaded from cache for coordinates (%.4f, %.4f)", lat, long)
            return nil
        }
    }
    
    var respData response
    
    params := fmt.Sprintf(
        "latitude=%f&longitude=%f&current=temperature_2m",
        lat,
        long,
    )
    
    url := fmt.Sprintf("%s?%s", apiURL, params)
    wi.l.Debugf("url was generated success - %s", url)
    
    resp, err := http.Get(url)
    if err != nil {
        wi.l.Error("can't get weather data", err)
        customErr := errors.New("can't get weather data from openmeteo")
        return errors.Join(customErr, err)
    }
    defer func() {
        if err := resp.Body.Close(); err != nil {
            wi.l.Error("can't close body", err)
        }
    }()
    
    data, err := io.ReadAll(resp.Body)
    if err != nil {
        wi.l.Error("can't read data from body", err)
        customErr := errors.New("can't read data from response")
        return errors.Join(customErr, err)
    }
    
    wi.l.Debugf("data was readed successfully size - %d", len(data))
    
    if err := json.Unmarshal(data, &respData); err != nil {
        wi.l.Error("can't unmarshal json data", err)
        customErr := errors.New("can't unmarshal data from response")
        return errors.Join(customErr, err)
    }
    
    wi.c = respData.Curr
    wi.isLoaded = true
    
    if err := wi.cache.Set(cacheKey, float64(wi.c.Temp)); err != nil {
        wi.l.Errorf("Failed to save to cache: %v", err)
    } else {
        wi.l.Debugf("Data saved to cache for coordinates (%.4f, %.4f)", lat, long)
    }
    
    return nil
}

func (wi *WeatherInfo) GetTemperature(lat, long float64) (models.TempInfo, error) {
    err := wi.getWeatherInfo(lat, long)
    return models.TempInfo{
        Temp: wi.c.Temp,
    }, err
}

func (wi *WeatherInfo) RefreshTemperature(lat, long float64) (models.TempInfo, error) {
    cacheKey := wi.getCacheKey(lat, long)
    wi.cache.Delete(cacheKey)
    wi.isLoaded = false
    err := wi.getWeatherInfo(lat, long)
    return models.TempInfo{
        Temp: wi.c.Temp,
    }, err
}

func (wi *WeatherInfo) CacheStats() int {
    return wi.cache.Size()
}
