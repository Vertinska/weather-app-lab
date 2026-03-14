package config

// LocationStorage - интерфейс для хранения местоположения
type LocationStorage interface {
    GetLocation() (float64, float64, error)
    SetLocation(lat, lon float64) error
}

// Location - структура с координатами
type Location struct {
    Latitude  float64
    Longitude float64
}
