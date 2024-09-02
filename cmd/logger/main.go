package logger

import "github.com/charmbracelet/log"

func LogInfo(message string) {
	log.Printf("INFO", message)
}

func LogWarning(message string) {
	log.Printf("WARN - %v", message)
}

func LogError(message string) {
	log.Printf("ERROR - %v", message)
}
