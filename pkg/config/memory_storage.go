package config

type MemoryStorage struct {
    latitude  float64
    longitude float64
}

func NewMemoryStorage() *MemoryStorage {
    // По умолчанию Гродно
    return &MemoryStorage{
        latitude:  53.6688,
        longitude: 23.8223,
    }
}

func (m *MemoryStorage) GetLocation() (float64, float64, error) {
    return m.latitude, m.longitude, nil
}

func (m *MemoryStorage) SetLocation(lat, lon float64) error {
    m.latitude = lat
    m.longitude = lon
    return nil
}
