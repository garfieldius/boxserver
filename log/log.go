package log

import (
	"fmt"
	"time"
	"os"
)

const timeLayout = "2006-01-02 15:04:05.999"

func stdOut(message string, args ...interface {}) {
	t := time.Now().Format(timeLayout)
	fmt.Printf(t + " " + message + "\n", args...);
}

func stdErr(message string, args ...interface {}) {
	t := time.Now().Format(timeLayout)
	fmt.Fprintf(os.Stderr, t + " " + message + "\n", args...)
}

func Debug(message string, args ...interface {}) {
	stdOut("[DEBUG] " + message, args...)
}

func Notice(message string, args... interface {}) {
	stdOut("[NOTICE] " + message, args...)
}

func Warn(message string, args... interface {}) {
	stdErr("[WARNING] " + message, args...)
}

func Error(message string, args... interface {}) {
	stdErr("[ERROR] " + message, args...)
}

func Critical(message string, args... interface {}) {
	stdErr("[CRITICAL] " + message, args...)
	os.Exit(1)
}
