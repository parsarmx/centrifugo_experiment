package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	viper.SetConfigFile(path)

	// Set default values
	viper.SetDefault("environment", "development")

	viper.AutomaticEnv()
	// viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config values
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func GetDSN(cfg *DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)
}

func NewLoggerEncoderConfig(cfg *LoggerConfig) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        cfg.TimeKey,
		LevelKey:       cfg.LevelKey,
		NameKey:        cfg.NameKey,
		CallerKey:      cfg.CallerKey,
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     cfg.MessageKey,
		StacktraceKey:  cfg.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// func GetRabbitMQUrl(cfg *RabbitConfig) string {
// 	return fmt.Sprintf("amqp://%s:%s@%s:%s/",
// 		cfg.User,
// 		cfg.Password,
// 		cfg.Host,
// 		cfg.Port,
// 	)
// }
