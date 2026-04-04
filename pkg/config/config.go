package config

import (
    "io"
    "os"
    "gopkg.in/yaml.v3"
)

type ConfigFile struct {
    C Config `yaml:"service"`
}

type Provider struct {
    Type     string `yaml:"type"`
    CacheTTL int    `yaml:"cache_ttl"`
}

type Location struct {
    Lat  float64 `yaml:"lat"`
    Long float64 `yaml:"long"`
}

type Config struct {
    P Provider `yaml:"provider"`
    L Location `yaml:"location"`
}

func Parse(r io.Reader) (Config, error) {
    var c ConfigFile
    if err := yaml.NewDecoder(r).Decode(&c); err != nil {
        return Config{}, err
    }
    return c.C, nil
}

func LoadConfigFromFile(path string) (Config, error) {
    r, err := os.Open(path)
    if err != nil {
        return Config{}, err
    }
    defer r.Close()
    return Parse(r)
}
