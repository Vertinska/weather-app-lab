package logger

import (
    "fmt"
    "log"
    "os"
)

type MyLogger struct {
    infoLogger  *log.Logger
    debugLogger *log.Logger
    errorLogger *log.Logger
    debugMode   bool
}

func New(debugMode bool) *MyLogger {
    return &MyLogger{
        infoLogger:  log.New(os.Stdout, "📢 INFO: ", log.Ldate|log.Ltime),
        debugLogger: log.New(os.Stdout, "🔍 DEBUG: ", log.Ldate|log.Ltime),
        errorLogger: log.New(os.Stderr, "❌ ERROR: ", log.Ldate|log.Ltime),
        debugMode:   debugMode,
    }
}

func (l *MyLogger) Info(msg string) {
    l.infoLogger.Println(msg)
}

func (l *MyLogger) Debug(msg string) {
    if l.debugMode {
        l.debugLogger.Println(msg)
    }
}

func (l *MyLogger) Error(msg string) {
    l.errorLogger.Println(msg)
}

func (l *MyLogger) Infof(format string, args ...interface{}) {
    l.Info(fmt.Sprintf(format, args...))
}

func (l *MyLogger) Debugf(format string, args ...interface{}) {
    l.Debug(fmt.Sprintf(format, args...))
}

func (l *MyLogger) Errorf(format string, args ...interface{}) {
    l.Error(fmt.Sprintf(format, args...))
}
