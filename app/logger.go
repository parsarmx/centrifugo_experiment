package application

import (
	"fmt"

	"golang_template/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func (a *application) InitLogger() (*zap.Logger, error) {
	encoderCfg := config.NewLoggerEncoderConfig(&a.config.Logger)
	fileEncoder := zapcore.NewJSONEncoder(encoderCfg)
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   a.config.Logger.Lumberjack.Filename,
		MaxSize:    a.config.Logger.Lumberjack.MaxSize,
		MaxAge:     a.config.Logger.Lumberjack.MaxAge,
		MaxBackups: a.config.Logger.Lumberjack.MaxBackups,
		LocalTime:  a.config.Logger.Lumberjack.LocalTime,
		Compress:   a.config.Logger.Lumberjack.Compress,
	})

	logLevel, err := zapcore.ParseLevel(a.config.Logger.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse logger level:%w", err)
	}

	core := zapcore.NewCore(fileEncoder, file, logLevel)
	logger := zap.New(core)

	defer logger.Sync()

	return logger, nil
}
