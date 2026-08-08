package logger

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Logger struct {
	Verbose bool
}

func New(verbose bool) *Logger {
	return &Logger{Verbose: verbose}
}

func (l *Logger) Info(msg string) {
	color.Set(color.FgBlue)
	fmt.Printf("[%s] ℹ️ %s\n", time.Now().Format("15:04:05"), msg)
	color.Unset()
}

func (l *Logger) Success(msg string) {
	color.Set(color.FgGreen)
	fmt.Printf("[%s] ✅ %s\n", time.Now().Format("15:04:05"), msg)
	color.Unset()
}

func (l *Logger) Warn(msg string) {
	color.Set(color.FgYellow)
	fmt.Printf("[%s] ⚠️ %s\n", time.Now().Format("15:04:05"), msg)
	color.Unset()
}

func (l *Logger) Error(msg string) {
	color.Set(color.FgRed)
	fmt.Printf("[%s] ❌ %s\n", time.Now().Format("15:04:05"), msg)
	color.Unset()
}

func (l *Logger) Debug(msg string) {
	if l.Verbose {
		color.Set(color.FgMagenta)
		fmt.Printf("[%s] 🔍 %s\n", time.Now().Format("15:04:05"), msg)
		color.Unset()
	}
}

