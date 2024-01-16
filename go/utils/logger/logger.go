package logger

import (
	"fmt"
	"net/http"
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

func (l Logger) INFO(message string, args ...interface{}) {

	message = fmt.Sprintf("%s[INFO] %s%s\n", BOLD_GREEN, message, RESET)
	fmt.Printf(message, args...)

}

func (l Logger) DEBUG(message string, args ...interface{}) {

	if l.Environment != "PRODUCTION" {
		message = fmt.Sprintf("%s[DEBUG] %s%s\n", LIGHT_BLUE, message, RESET)
		fmt.Printf(message, args...)
	}

}

func (l Logger) DEBUG_WARN(message string, args ...interface{}) {

	if l.Environment != "PRODUCTION" {
		message = fmt.Sprintf("%s[DEBUG_WARN] %s%s\n", LIGHT_BLUE, message, RESET)
		fmt.Printf(message, args...)
	}

}

func (l Logger) WARN(message string, args ...interface{}) {

	message = fmt.Sprintf("%s[WARN] %s%s\n", BOLD_ORANGE, message, RESET)
	fmt.Printf(message, args...)

}

func (l Logger) ERROR(message string, args ...interface{}) {

	message = fmt.Sprintf("%s[ERROR] %s%s\n", BOLD_RED, message, RESET)
	fmt.Printf(message, args...)

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