package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	DEBUG = iota // Ярко-голубой
	INFO
	WARN  // Желтый
	ERROR // Красный
	FATAL // Красный BG
)

// Colors
const (
	BLACK          = "\033[30m" // Черный
	RED            = "\033[31m" // Красный
	GREEN          = "\033[32m" // Зеленый
	YELLOW         = "\033[33m" // Желтый
	BLUE           = "\033[34m" // Синий
	MAGENTA        = "\033[35m" // Пурпурный
	CYAN           = "\033[36m" // Голубой
	WHITE          = "\033[37m" // Белый
	BRIGHT_BLACK   = "\033[90m" // Серый
	BRIGHT_RED     = "\033[91m" // Ярко-красный
	BRIGHT_GREEN   = "\033[92m" // Ярко-зеленый
	BRIGHT_YELLOW  = "\033[93m" // Ярко-желтый
	BRIGHT_BLUE    = "\033[94m" // Ярко-синий
	BRIGHT_MAGENTA = "\033[95m" // Ярко-пурпурный
	BRIGHT_CYAN    = "\033[96m" // Ярко-голубой
	BRIGHT_WHITE   = "\033[97m" // Ярко-белый
)

// Background
const (
	BG_BLACK          = "\033[40m"  // Черный
	BG_RED            = "\033[41m"  // Красный
	BG_GREEN          = "\033[42m"  // Зеленый
	BG_YELLOW         = "\033[43m"  // Желтый
	BG_BLUE           = "\033[44m"  // Синий
	BG_MAGENTA        = "\033[45m"  // Пурпурный
	BG_CYAN           = "\033[46m"  // Голубой
	BG_WHITE          = "\033[47m"  // Белый
	BG_BRIGHT_BLACK   = "\033[100m" // Серый
	BG_BRIGHT_RED     = "\033[101m" // Ярко-красный
	BG_BRIGHT_GREEN   = "\033[102m" // Ярко-зеленый
	BG_BRIGHT_YELLOW  = "\033[103m" // Ярко-желтый
	BG_BRIGHT_BLUE    = "\033[104m" // Ярко-синий
	BG_BRIGHT_MAGENTA = "\033[105m" // Ярко-пурпурный
	BG_BRIGHT_CYAN    = "\033[106m" // Ярко-голубой
	BG_BRIGHT_WHITE   = "\033[107m" // Ярко-белый
)

// Attributes
const (
	RESET         = "\033[0m" // Сброс всех настроек
	BOLD          = "\033[1m" // Жирный текст
	DIM           = "\033[2m" // Блеклый текст
	ITALIC        = "\033[3m" // Курсив
	UNDERLINE     = "\033[4m" // Подчеркнутый
	BLINK         = "\033[5m" // Мигание
	INVERT        = "\033[7m" // Инверсия цветов
	HIDDEN        = "\033[8m" // Скрытый текст
	STRIKETHROUGH = "\033[9m" // Зачеркнутый текст
)

var levelMap = map[int]string{
	0: "DEBUG",
	1: "INFO",
	2: "WARN",
	3: "ERROR",
	4: "FATAL",
}

type CustomLogger struct {
	logger *log.Logger
	level  int
	mu     sync.Mutex
}

func NewLogger(level int, logger *log.Logger) *CustomLogger {
	return &CustomLogger{
		level:  level,
		logger: logger,
	}
}

func (l *CustomLogger) SetLevel(level int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if level < 5 && level >= 0 {
		l.level = level
	}
}

func (l *CustomLogger) logMessage(level int, color, format string, v ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if level < l.level {
		return
	}

	prefix := fmt.Sprintf("[%s%s\033[0m] | %s | ", color, levelMap[level], time.Now().Format("2021-02-03 15:14:32"))
	msg := fmt.Sprintf(format, v...)
	l.logger.Println(prefix, msg)
}

func (l *CustomLogger) Debug(format string, v ...interface{}) {
	l.logMessage(DEBUG, BRIGHT_CYAN, format, v...)
}

func (l *CustomLogger) Info(format string, v ...interface{}) {
	l.logMessage(INFO, WHITE, format, v...)
}

func (l *CustomLogger) Warn(format string, v ...interface{}) {
	l.logMessage(WARN, YELLOW, format, v...)
}

func (l *CustomLogger) Error(format string, v ...interface{}) {
	l.logMessage(ERROR, RED, format, v...)
}

func (l *CustomLogger) Fatal(format string, v ...interface{}) {
	l.logMessage(FATAL, BG_RED, format, v...)
	os.Exit(1)
}
