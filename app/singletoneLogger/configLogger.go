package singletoneLogger

import (
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"io"
	"os"
)

func init() {
	colorFunc = map[string]func(s string, a ...interface{}) string{
		"red":   color.New(color.FgRed).SprintfFunc(),
		"green": color.New(color.FgGreen).SprintfFunc(),
	}
}

// loggerData - структура, которая предоставляет данные для инициализации логгера.
type loggerData struct {
	Out      io.Writer                               // writer для логов
	BuffSize int                                     // максимальный размер каналов
	ErrColor func(s string, a ...interface{}) string // функция окрашивающая цвет для ошибок
	MsgColor func(s string, a ...interface{}) string // функция окрашивающая цвет для сообщений
}

var (
	errIncorrectValue = errors.New("Incorrect value")
	colorFunc         map[string]func(s string, a ...interface{}) string // мап для определения функции по имени
)

// newLoggerData - создает loggerData с значениями по умолчанию
func newLoggerData() *loggerData {
	return &loggerData{
		Out:      os.Stdout,
		BuffSize: 100,
		ErrColor: color.New(color.FgRed).SprintfFunc(),
		MsgColor: color.New(color.FgGreen).SprintfFunc(),
	}
}
