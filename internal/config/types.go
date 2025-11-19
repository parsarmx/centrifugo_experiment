package config

type Config struct {
	Server ServerConfig   `mapstructure:"server" validate:"required"`
	DB     DatabaseConfig `mapstructure:"db" validate:"required"`
	Logger LoggerConfig   `mapstructure:"logger" validate:"required"`
	// Rabbit RabbitConfig   `mapstructure:"rabbit" validate:"required"`
	Redis RedisConfig `mapstructure:"redis" validate:"required"`
	Auth  AuthConfig  `mapstructure:"auth" validate:"required"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host" validate:"required,hostname|ip"`
	Port     string `mapstructure:"port" validate:"required,number"`
	User     string `mapstructure:"user" validate:"required"`
	Password string `mapstructure:"password" validate:"required"`
	DBName   string `mapstructure:"dbname" validate:"required"`
	SSLMode  string `mapstructure:"sslmode" validate:"required,oneof=disable enable verify-full"`
	MaxConns int    `mapstructure:"max_conns" validate:"required,min=1"`
	MinConns int    `mapstructure:"min_conns" validate:"required,min=1"`
}

type ServerConfig struct {
	Port         string `mapstructure:"port" validate:"required,number"`
	Host         string `mapstructure:"host" validate:"required,hostname|ip"`
	Mode         string `mapstructure:"mode" validate:"required,oneof=development production testing"`
	ReadTimeout  int    `mapstructure:"read_timeout" validate:"required,min=1"`
	WriteTimeout int    `mapstructure:"write_timeout" validate:"required,min=1"`
}

type LoggerConfig struct {
	Level         string           `mapstructure:"level" validate:"required,oneof=debug info warn error dpanic panic fatal"`
	TimeKey       string           `mapstructure:"time_key" validate:"required"`
	LevelKey      string           `mapstructure:"level_key" validate:"required"`
	NameKey       string           `mapstructure:"name_key" validate:"required"`
	CallerKey     string           `mapstructure:"caller_key" validate:"required"`
	MessageKey    string           `mapstructure:"message_key" validate:"required"`
	StacktraceKey string           `mapstructure:"stacktrace_key" validate:"required"`
	Lumberjack    LumberjackConfig `mapstructure:"lumberjack" validate:"required"`
}

type LumberjackConfig struct {
	Filename   string `mapstructure:"filename" validate:"required"`
	MaxSize    int    `mapstructure:"max_size" validate:"required,number"`
	MaxAge     int    `mapstructure:"max_age" validate:"required,number"`
	MaxBackups int    `mapstructure:"max_backups" validate:"required,number"`
	LocalTime  bool   `mapstructure:"local_time"`
	Compress   bool   `mapstructure:"compress"`
}

// type RabbitConfig struct {
// 	Host     string `mapstructure:"host" validate:"required,hostname|ip"`
// 	Port     string `mapstructure:"port" validate:"required,number"`
// 	User     string `mapstructure:"user" validate:"required"`
// 	Password string `mapstructure:"password" validate:"required"`
// }

type AuthConfig struct {
	OTPCodeLength    int            `mapstructure:"otp_code_length" validate:"required"`
	OTPTTL           int            `mapstructure:"otp_ttl" validate:"required"` // seconds
	UnderDevelopment bool           `mapstructure:"under_development"`
	JWTSecret        string         `mapstructure:"jwt_secret" validate:"required"`
	TestUser         TestUserConfig `mapstructure:"test_user"`
	AccTokenExpTime  int            `mapstructure:"access_token_exp_time" validate:"required"` // minutes
}

type TestUserConfig struct {
	PhoneNumber string `mapstructure:"phone_number" validate:"required,e164"`
	OTPCode     string `mapstructure:"otp_code" validate:"required,numeric"`
}
type RedisConfig struct {
	Host         string        `mapstructure:"host" validate:"required"`
	Port         int           `mapstructure:"port" validate:"required,number"`
	Password     string        `mapstructure:"password"`
	DB           int           `mapstructure:"db"`
	MaxRetries   int           `mapstructure:"max_retries" validate:"required,min=1"`
	PoolSize     int           `mapstructure:"pool_size" validate:"required,min=1"`
	MinIdleConns int           `mapstructure:"min_idle_conns" validate:"required,min=1"`
	Timeouts     RedisTimeouts `mapstructure:"timeouts" validate:"required"`
}

type RedisTimeouts struct {
	Dial  int `mapstructure:"dial" validate:"required,min=1"`
	Read  int `mapstructure:"read" validate:"required,min=1"`
	Write int `mapstructure:"write" validate:"required,min=1"`
	Idle  int `mapstructure:"idle" validate:"required,min=1"`
}
