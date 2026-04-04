package cache

import (
    "encoding/json"
    "os"
    "path/filepath"
    "sync"
    "time"
)

type FileCache struct {
    cacheDir string
    mu       sync.RWMutex
    ttl      time.Duration
}

type cacheData struct {
    Value      interface{} `json:"value"`
    Expiration time.Time   `json:"expiration"`
}

func NewFileCache(cacheDir string, ttl time.Duration) *FileCache {
    // Создаем директорию для кеша, если её нет
    os.MkdirAll(cacheDir, 0755)
    
    return &FileCache{
        cacheDir: cacheDir,
        ttl:      ttl,
    }
}

func (c *FileCache) getFilePath(key string) string {
    return filepath.Join(c.cacheDir, key+".json")
}

func (c *FileCache) Set(key string, value interface{}) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    data := cacheData{
        Value:      value,
        Expiration: time.Now().Add(c.ttl),
    }
    
    filePath := c.getFilePath(key)
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }
    
    return os.WriteFile(filePath, jsonData, 0644)
}

func (c *FileCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    filePath := c.getFilePath(key)
    
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, false
    }
    
    var cacheData cacheData
    if err := json.Unmarshal(data, &cacheData); err != nil {
        return nil, false
    }
    
    // Проверяем, не истек ли срок
    if time.Now().After(cacheData.Expiration) {
        os.Remove(filePath)
        return nil, false
    }
    
    return cacheData.Value, true
}

func (c *FileCache) Delete(key string) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    filePath := c.getFilePath(key)
    return os.Remove(filePath)
}

func (c *FileCache) Clear() error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    return os.RemoveAll(c.cacheDir)
}

func (c *FileCache) Size() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    files, err := os.ReadDir(c.cacheDir)
    if err != nil {
        return 0
    }
    return len(files)
}
