package logger

import (
    "fmt"
    "time"
)

var (
    INFO  = "INFO"
    DEBUG = "DEBUG"
    ERROR = "ERROR"
)

type logger struct{}

func New() *logger {
    return &logger{}
}

func (l *logger) Info(msg string) {
    fmt.Println(l.msg(INFO, msg))
}

func (l *logger) Debug(msg string) {
    fmt.Println(l.msg(DEBUG, msg))
}

func (l *logger) Error(msg string, err error) {
    fmt.Println(l.msg(ERROR, msg+" err - "+err.Error()))
}

func (l *logger) msg(level string, msg string) string {
    timeStr := time.Now().Format(time.RFC3339)
    return fmt.Sprintf(
        "[%s] %s, message - %s",
        level,
        timeStr,
        msg,
    )
}
