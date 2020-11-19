// Package log is a very lightweight wrapper around the fatih/color and zerolog packages for log output
package log

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	blue   = color.New(color.FgBlue, color.Bold)
	yellow = color.New(color.FgYellow, color.Bold)
	green  = color.New(color.FgGreen, color.Bold)
	red    = color.New(color.FgRed, color.Bold)
)

// Initialize Initializes logging configuration
func Initialize(logLevel string) {
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}

	consoleWriter.FormatTimestamp = func(i interface{}) string {
		return ""
	}
	consoleWriter.FormatLevel = func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "debug":
				l = prefix(blue, ll)
			case "info":
				l = prefix(green, ll)
			case "warn":
				l = prefix(yellow, "warning")
			case "error":
				l = prefix(red, ll)
			case "fatal":
				l = prefix(red, ll)
			}
		}

		return l
	}

	multi := zerolog.MultiLevelWriter(consoleWriter)

	log.Logger = zerolog.New(multi)

	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		Fatalf("Unknown log level %s", logLevel)
	}

	Debugln("Logger has been configured")
}

// Level returns the current level of the logger
func Level() string {
	return zerolog.GlobalLevel().String()
}

func prefix(c *color.Color, p string) string {
	return c.SprintfFunc()(p + ":")
}

// Debugln prints a debug level msg
func Debugln(a ...interface{}) {
	log.Debug().Msg(fmt.Sprint(a...))
}

// Debugf prints a debug level msg
func Debugf(format string, a ...interface{}) {
	log.Debug().Msgf(format, a...)
}

// Infof prints an information level msg
func Infof(format string, a ...interface{}) {
	log.Info().Msgf(format, a...)
}

// Infoln prints an information level msg
func Infoln(a ...interface{}) {
	log.Info().Msg(fmt.Sprint(a...))
}

// Warningf prints a warning level msg
func Warningf(format string, a ...interface{}) {
	log.Warn().Msgf(format, a...)
}

// Warningln prints a warning level msg
func Warningln(a ...interface{}) {
	log.Warn().Msg(fmt.Sprint(a...))
}

// Errorf prints an error level msg
func Errorf(format string, a ...interface{}) {
	log.Error().Msgf(format, a...)
}

// Errorln prints an error level msg
func Errorln(a ...interface{}) {
	log.Error().Msg(fmt.Sprint(a...))
}

// Fatalf prints a fatal level msg
func Fatalf(format string, a ...interface{}) {
	log.Fatal().Msgf(format, a...)
}

//Fatalln prints a fatal level msg
func Fatalln(a ...interface{}) {
	log.Fatal().Msg(fmt.Sprint(a...))
}
