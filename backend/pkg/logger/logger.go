package logger

import (
	"encoding/json"
	"log"
)

type Logger struct {
	prefix string
}

func New(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) LogStructured(level, message string, fields map[string]interface{}) {
	fieldsJSON, _ := json.Marshal(fields)
	log.Printf("[%s][%s] %s | %s", l.prefix, level, message, string(fieldsJSON))
}

func (l *Logger) Debug(message string, fields map[string]interface{}) {
	l.LogStructured("DEBUG", message, fields)
}

func (l *Logger) Info(message string, fields map[string]interface{}) {
	l.LogStructured("INFO", message, fields)
}

func (l *Logger) Error(message string, fields map[string]interface{}) {
	l.LogStructured("ERROR", message, fields)
}