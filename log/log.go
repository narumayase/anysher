package log

import (
	"github.com/rs/zerolog"
	"os"
	"strings"
)

// SetLogLevel sets the global zerolog log level based on the provided string.
// If the logLevel string is not recognized, it defaults to zerolog.InfoLevel.
func SetLogLevel() {
	levels := map[string]zerolog.Level{
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"fatal": zerolog.FatalLevel,
		"panic": zerolog.PanicLevel,
	}
	levelEnv := strings.ToLower(os.Getenv("LOG_LEVEL"))

	level, ok := levels[levelEnv]
	if !ok {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}
