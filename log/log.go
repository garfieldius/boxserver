package log

import (
	"fmt"
	"time"
	"os"
)

const timeLayout string = "2006-01-02 15:04:05"

func now() string {
	return time.Now().Format(timeLayout)
}

func stdOut(message string, args ...interface {}) {
	fmt.Printf(now() + " " + message + "\n", args...);
}

func stdErr(message string, args ...interface {}) {
	fmt.Fprintf(os.Stderr, now() + " " + message + "\n", args...)
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
