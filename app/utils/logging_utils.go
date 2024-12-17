package utils

import "log"

var logColors map[string]string = map[string]string{

	"reset":   "\033[0m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"gray":    "\033[37m",
	"white":   "\033[97m",
}

func Log(msg ...any) {
	logColor("white", msg...)
}

func LogSuccess(msg ...any) {
	logColor("green", msg...)
}

func LogWarning(msg ...any) {
	logColor("yellow", msg...)
}

func LogError(msg ...any) {
	logColor("red", msg...)
}

func logColor(color string, msg ...any) {
	colorCode := logColors[color]
	resetCode := logColors["reset"]
	if colorCode == "" {
		resetCode = ""
	}
	msg = append([]any{colorCode}, msg...)
	msg = append(msg, resetCode)
	log.Println(msg...)
}
