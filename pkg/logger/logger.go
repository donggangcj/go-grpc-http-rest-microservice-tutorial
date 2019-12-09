package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

var (
	// Log is global logger
	Log *zap.Logger

	// timeFormat is custom Time format
	customTimeFormat string

	// onceInit guarantee initialize logger only once
	onceInit sync.Once
)

// customTimeEncode encode time to our format
// This example how we can customize zap default functionality
func customTimeEncode(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(customTimeFormat))
}

//Init initialize log by input parameters
// lv1 - global log level: Debug(-1),Info(0),Warn(1),Error(2),DPanic(4),Fatal(5)
func Init(lvl int, timeFormat string) error {
	var err error
	onceInit.Do(func() {
		globalLevle := zapcore.Level(lvl)

		highPriority := zap.LevelEnablerFunc(func(lv1 zapcore.Level) bool {
			return lv1 >= zapcore.ErrorLevel
		})

		lowPriority := zap.LevelEnablerFunc(func(lv1 zapcore.Level) bool {
			return lv1 >= globalLevle && lv1 < zapcore.ErrorLevel
		})

		consoleInfos := zapcore.Lock(os.Stdout)
		consoleErrors := zapcore.Lock(os.Stderr)

		// Configure console output
		var useCustomTimeFormat bool
		ecfg := zap.NewProductionEncoderConfig()
		if len(timeFormat) > 0 {
			customTimeFormat = timeFormat
			ecfg.EncodeTime = customTimeEncode
			useCustomTimeFormat = true
		}
		consoleEncoder := zapcore.NewJSONEncoder(ecfg)

		// Join the outputs,encoders,and level-handling functions into zapcore
		core := zapcore.NewTee(zapcore.NewCore(consoleEncoder, consoleErrors, highPriority), zapcore.NewCore(consoleEncoder, consoleInfos, lowPriority))

		// From a zapcore.Core, it's easy  to construct a Logger
		Log = zap.New(core)

		if !useCustomTimeFormat {
			Log.Warn("time format for logger is note provided - use zap default")
		}
	})
	return err

}
