package logger

type Logger interface {
    Info(msg string)
    Debug(msg string)
    Error(msg string)
    Infof(format string, args ...interface{})
    Debugf(format string, args ...interface{})
    Errorf(format string, args ...interface{})
}
