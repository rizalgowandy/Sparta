package sparta

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	colorRed     = 31
	colorYellow  = 33
	colorBlue    = 34
	colorMagenta = 35

	colorBold = 1
)

func newRSLogger(logLevel zerolog.Level, outputFormat string, noColor bool) (*zerolog.Logger, error) {
	var loggerWriter io.Writer
	switch outputFormat {
	case "text", "txt":
		consoleWriter := zerolog.ConsoleWriter{
			Out:        colorable.NewColorableStdout(),
			TimeFormat: time.RFC822,
		}
		consoleWriter.FormatLevel = func(i interface{}) string {
			var l string
			if ll, ok := i.(string); ok {
				switch ll {
				case "trace":
					l = colorize("|TRACE|", colorMagenta, noColor)
				case "debug":
					l = colorize("|DEBUG|", colorYellow, noColor)
				case "info":
					l = colorize("|INFO |", colorBlue, noColor)
				case "warn":
					l = colorize("|WARN |", colorYellow, noColor)
				case "error":
					l = colorize(colorize("|ERROR|", colorRed, noColor), colorBold, noColor)
				case "fatal":
					l = colorize(colorize("|FATAL|", colorRed, noColor), colorBold, noColor)
				case "panic":
					l = colorize(colorize("|PANIC|", colorRed, noColor), colorBold, noColor)
				default:
					l = colorize("???", colorBold, noColor)
				}
			} else {
				if i == nil {
					l = colorize("???", colorBold, noColor)
				} else {
					l = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
				}
			}
			return l
		}
		consoleWriter.FormatMessage = func(i interface{}) string {
			// 48 is the same as the dividerLength value
			return fmt.Sprintf("%-48s", i)
		}
		consoleWriter.FormatFieldName = func(i interface{}) string {
			return colorize(fmt.Sprintf("%s=", i), colorBold, noColor)
		}
		consoleWriter.FormatFieldValue = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		loggerWriter = &consoleWriter
	default:
		loggerWriter = os.Stdout
	}
	// Set it up and return it...
	rsLogger := zerolog.New(loggerWriter).With().Timestamp().Logger().Level(logLevel)
	updateLoggerGlobals()
	return &rsLogger, nil
}

// NewLogger returns a new zerolog.Logger instance. It is the caller's responsibility
// to set the formatter if needed.
func NewLogger(level string) (*zerolog.Logger, error) {
	return NewLoggerForOutput(level, "", false)
}

// NewLoggerForOutput returns a new zerolog
func NewLoggerForOutput(userLevel string, outputType string, disableColors bool) (*zerolog.Logger, error) {
	// If there is an environment override, use that
	envLogLevel := os.Getenv(envVarLogLevel)
	if envLogLevel != "" {
		userLevel = envLogLevel
	}
	logLevel, err := zerolog.ParseLevel(userLevel)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse logLevel: %s", userLevel)
	}
	return newRSLogger(logLevel, outputType, disableColors)
}
