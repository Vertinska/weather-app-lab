package logger

type Logger interface {
    Info(string)
    Debug(string)
    Error(string, error)
    Infof(string, ...interface{})
    Debugf(string, ...interface{})
    Errorf(string, ...interface{})
}
