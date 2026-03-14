package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type FileStorage struct {
    filePath string
}

type locationData struct {
    Latitude  float64 `json:"latitude"`
    Longitude float64 `json:"longitude"`
}

func NewFileStorage(filePath string) *FileStorage {
    // Если путь не указан, используем стандартный
    if filePath == "" {
        homeDir, _ := os.UserHomeDir()
        filePath = filepath.Join(homeDir, ".weather-app", "location.json")
    }
    
    // Создаем директорию, если её нет
    os.MkdirAll(filepath.Dir(filePath), 0755)
    
    return &FileStorage{
        filePath: filePath,
    }
}

func (f *FileStorage) GetLocation() (float64, float64, error) {
    data, err := os.ReadFile(f.filePath)
    if err != nil {
        if os.IsNotExist(err) {
            // Если файла нет, возвращаем координаты Гродно по умолчанию
            return 53.6688, 23.8223, nil
        }
        return 0, 0, err
    }
    
    var loc locationData
    if err := json.Unmarshal(data, &loc); err != nil {
        return 0, 0, err
    }
    
    return loc.Latitude, loc.Longitude, nil
}

func (f *FileStorage) SetLocation(lat, lon float64) error {
    loc := locationData{
        Latitude:  lat,
        Longitude: lon,
    }
    
    data, err := json.MarshalIndent(loc, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(f.filePath, data, 0644)
}
