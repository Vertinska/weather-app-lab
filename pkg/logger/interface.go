package logger

type Logger interface {
    Info(string)
    Debug(string)
    Error(string, error)
}
