package config

import (
    "encoding/json"
    "os"
)

type Config struct {
    StorageType string `json:"storage_type"`
    FilePath    string `json:"file_path,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
    // Если файла нет, создаем конфиг по умолчанию
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return &Config{
                StorageType: "file",
                FilePath:    "",
            }, nil
        }
        return nil, err
    }
    
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    return &config, nil
}

func (c *Config) Save(path string) error {
    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }
    
    return os.WriteFile(path, data, 0644)
}
