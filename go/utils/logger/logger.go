package logger

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	BOLD_RED    = "\033[1;31m"       // error
	BOLD_ORANGE = "\033[1;38;5;208m" // warn
	BOLD_GREEN  = "\033[1;32m"       // info
	LIGHT_BLUE  = "\033[36m"         // debug
	

	RESET = "\033[0m"
)

type Logger struct {
	Environment string
}

type loggingHandlerFunc func(w http.ResponseWriter, r *http.Request, l *Logger, requestID string)

func (l *Logger) Middleware(next loggingHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		next(w, r, l, requestID)
	}
}

func (l *Logger) SetEnvironment(environment string) {
	l.Environment = environment
}


func (l *Logger) INFO(message string, args ...interface{}) {
    _, file, line, _ := runtime.Caller(1)
    padding := calculatePadding(file, line)
    logMessage := fmt.Sprintf("%s[%s:%d] %s [INFO] %s %s\n", BOLD_GREEN, file, line, padding, message, RESET)
    fmt.Printf(logMessage, args...)
}

func (l *Logger) DEBUG(message string, args ...interface{}) {
    if l.Environment != "PRODUCTION" {
        _, file, line, _ := runtime.Caller(1)
        padding := calculatePadding(file, line)
        logMessage := fmt.Sprintf("%s[%s:%d] %s [DEBUG] %s %s\n", LIGHT_BLUE, file, line, padding, message, RESET)
        fmt.Printf(logMessage, args...)
    }
}

func (l *Logger) DEBUG_WARN(message string, args ...interface{}) {
	if l.Environment != "PRODUCTION" {
		_, file, line, _ := runtime.Caller(1)
		padding := calculatePadding(file, line)
		logMessage := fmt.Sprintf("%s[%s:%d] %s [DEBUG_WARN] %s %s\n", BOLD_ORANGE, file, line, padding, message, RESET)
		fmt.Printf(logMessage, args...)
	}
}

func (l *Logger) WARN(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	padding := calculatePadding(file, line)
	logMessage := fmt.Sprintf("%s[%s:%d] %s [WARN] %s %s\n", BOLD_ORANGE, file, line, padding, message, RESET)
	fmt.Printf(logMessage, args...)
}

func (l *Logger) ERROR(message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	padding := calculatePadding(file, line)
	logMessage := fmt.Sprintf("%s[%s:%d] %s [ERROR] %s %s\n", BOLD_RED, file, line, padding, message, RESET)
	fmt.Printf(logMessage, args...)
}

func calculatePadding(file string, line int) string {
	paddingLength :=  75 - len(fmt.Sprintf("[%s:%d]", file, line))
	if paddingLength < 0 {
		paddingLength = 0
	}
	return strings.Repeat("_", paddingLength)
}

func ( Logger ) STARTTIME( ) time.Time {
	return time.Now()
}
func ( *Logger ) ENDTIME( startTime time.Time, formatString string, v ...interface{} ) {
	elapsed := time.Since(startTime).Seconds()
	elapsedTimeString := fmt.Sprintf("Time elapsed: %f", elapsed)

	if elapsed <= 0.5 {
		return
	}

	if elapsed > 10.0 {		
		formatString = fmt.Sprintf("%s[DEBUGWARNING] %s COMPLETED This took more than 10 seconds. %s%s\n", BOLD_ORANGE, formatString, elapsedTimeString, RESET)

	} else if elapsed > 0.5 {
		formatString = fmt.Sprintf("%s[DEBUGWARNING] %s COMPLETED This took more than 1/2 a second. %s%s\n", BOLD_ORANGE, formatString, elapsedTimeString, RESET)
	
	} 
	fmt.Printf( formatString, v... )
}