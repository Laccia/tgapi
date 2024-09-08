package pkg

import (
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog"
)

func newConsoleWriter(format LogFormat, dest io.Writer) io.Writer {
	w := zerolog.NewConsoleWriter()
	// устанавливаем значения вывода в пришедший из newLogWriter dest
	w.Out = dest

	switch format {
	case LogFormatCEFGPN:
		w.FormatTimestamp = func(interface{}) string {
			return fmt.Sprintf("end=%d", time.Now().Unix())
		}
	case LogFormatDefault:
		w.FormatTimestamp = func(i interface{}) string {
			return time.Now().UTC().Format(timeFormat)
		}
	default:
		w.FormatTimestamp = func(i interface{}) string {
			return time.Now().UTC().Format(timeFormat)
		}
	}

	return w
}
