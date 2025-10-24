package logging

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logLevel = new(slog.LevelVar)
)

// Setup slog defaults
func Setup(logfp string) io.Writer {
	multiw := io.MultiWriter(&lumberjack.Logger{
		Filename:   logfp,
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   false,
	}, os.Stdout)
	slog.SetDefault(slog.New(
		slog.NewTextHandler(multiw, &slog.HandlerOptions{Level: logLevel}),
	))
	return multiw
}

// Set loggin level from config
func SetLevel(debug bool) {
	if debug {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}
	slog.Info("Logging level set", "logging_level", logLevel)
}
