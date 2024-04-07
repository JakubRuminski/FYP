package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
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
	Verbose	    bool
	ClientID    string
	File        *os.File
}


func (l *Logger) InitRequestLogFile(clientID string) (file *os.File, ok bool) {
	filePath := fmt.Sprintf("/logs/%s.txt", clientID)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		l.ERROR("Error opening or creating log file: %s", err)
		return file, false
	}

	l.File = file
	return file, true
}


func (l *Logger) SetFlags(environment string, verbose bool, clientID string) {
	l.Environment = environment
	l.Verbose = verbose
	l.ClientID = clientID
}


func (l *Logger) INFO(message string, args ...interface{}) {
    l.writeToTerminal("INFO", BOLD_GREEN, message, args...)
	l.writeToFile("INFO", message, args...)
}

func (l *Logger) DEBUG(message string, args ...interface{}) {
    if l.Environment != "PRODUCTION" && l.Verbose {
        l.writeToTerminal("DEBUG", LIGHT_BLUE, message, args...)
    }
	l.writeToFile("DEBUG", message, args...)
}

func (l *Logger) DEBUG_WARN(message string, args ...interface{}) {
	if l.Environment != "PRODUCTION" && l.Verbose {
		l.writeToTerminal("DEBUG_WARN", BOLD_ORANGE, message, args...)
	}
	l.writeToFile("DEBUG_WARN", message, args...)
}

func (l *Logger) WARN(message string, args ...interface{}) {
	l.writeToTerminal("WARN", BOLD_ORANGE, message, args...)
	l.writeToFile("WARN", message, args...)
}

func (l *Logger) ERROR(message string, args ...interface{}) {
	l.writeToTerminal("ERROR", BOLD_RED, message, args...)
	l.writeToFile("ERROR", message, args...)
}






func (l *Logger) writeToTerminal(messageType, messageColor, message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	padding := calculatePadding(file, line)

	if ! l.Verbose && l.Environment == "UNIT_TESTING" {
		return
	}

	terminalMessage := fmt.Sprintf("%s[%s:%d] %s [%s] %s %s\n", messageColor, file, line, padding, messageType, message, RESET)
	log.Printf(terminalMessage, args...)
}


func (l *Logger) writeToFile(messageType, message string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	padding := calculatePadding(file, line)

	fileMessage     := fmt.Sprintf("[%s:%d] %s [%s] %s\n", file, line, padding, messageType, message)
	l.File.WriteString(fmt.Sprintf(fileMessage, args...))
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
    log.Printf( formatString, v... )
}