package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DefaultLogLevel = "debug"
)

var (
	Logger   *zap.Logger
	logLevel zap.AtomicLevel
)

func InitLogger(levelString string) {
	// Log to the console by default.
	logLevel = zap.NewAtomicLevel()
	setLogLevel(levelString)
	encoderCfg := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout), logLevel)
	Logger = zap.New(core, zap.AddCaller())
}

func ShutdownLogger() {
	_ = Logger.Sync()
}

func setLogLevel(level string) {
	if level != "" {
		parsedLevel, err := zapcore.ParseLevel(level)
		if err != nil {
			// Fallback to logging at the info level.
			fmt.Printf("Falling back to the info log level. You specified: %s.\n", level)
			logLevel.SetLevel(zapcore.InfoLevel)
		} else {
			logLevel.SetLevel(parsedLevel)
		}
	} else {
		logLevel.SetLevel(zapcore.DebugLevel)
	}
}
