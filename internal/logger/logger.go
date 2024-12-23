package logger

import (
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() error {
	env := viper.GetString("app.env")

	if env == "production" {
		lumberjackLogger := &lumberjack.Logger{
			Filename:   viper.GetString("log.file"),
			MaxSize:    viper.GetInt("log.max_size"),
			MaxBackups: viper.GetInt("log.max_backups"),
			MaxAge:     viper.GetInt("log.max_age"),
			Compress:   viper.GetBool("log.compress"),
		}
		enderConfig := zap.NewProductionEncoderConfig()
		encoder := zapcore.NewJSONEncoder(enderConfig)

		writeSyncer := zap.CombineWriteSyncers(
			zapcore.AddSync(lumberjackLogger),
			zapcore.AddSync(os.Stdout),
		)

		core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)

		Log = zap.New(core)
	} else {
		var err error
		Log, err = zap.NewDevelopment()
		if err != nil {
			return fmt.Errorf("failed to create development logger: %v", err)
		}
	}
	Log.Info("Init logger complete.")
	return nil
}

func AddRequestIDToLogger(requestId string) *zap.Logger {
	return Log.With(zap.String("X-Request-ID", requestId))
}

func SyncLogger() {
	Log.Info("Flush logger.")
	if err := Log.Sync(); err != nil {
		fmt.Println("Error syncing logger: ", err)
	}
}
