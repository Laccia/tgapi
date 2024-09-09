package pkg

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Config struct {
	Level       LogLevel
	Format      LogFormat
	Destination LogDestination
}

// parseLevel парсинг уровня логгирования из пришедшей структуры конфигов
func parseLevel(level LogLevel) (zerolog.Level, error) {
	lvl, err := zerolog.ParseLevel(string(level))
	if err != nil {
		return 0, fmt.Errorf("zerolog.ParseLevel failed: %w", err)
	}

	return lvl, nil
}

// newLogWriter создает экзепляр io.Writer с определением назначения вывода
func newLogWriter(dest LogDestination, format LogFormat) io.Writer {

	// указание, куда писать (по дефолту stdout)
	switch dest {
	case LogDestinationConsoleOut:
		return newConsoleWriter(format, os.Stdout)
	case LogDestinationConsoleErr:
		return newConsoleWriter(format, os.Stderr)
	default:
		return newConsoleWriter(format, os.Stdout)
	}
}

// NewLogger конструктор для инициализации логгера
func NewLogger(config Config) zerolog.Logger {
	// парсинг уровня
	level, err := parseLevel(config.Level)
	if err != nil {
		panic(err)
	}

	// создание writer с передачей формата и назначения (куда выводить)
	writer := newLogWriter(config.Destination, config.Format)

	logger := zerolog.
		New(writer).
		Level(level)
	return logger

}
