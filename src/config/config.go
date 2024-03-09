package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Logger   LoggerConfig
	Otp      OtpConfig
	JWT      JWTConfig
	Cors     CorseConfig
}
type CorseConfig struct {
	AllowOrigins string
}
type ServerConfig struct {
	Port    int
	RunMode string
}

type JWTConfig struct {
	Secret                     string
	RefreshSecret              string
	AccessTokenExpireDuration  time.Duration
	RefreshTokenExpireDuration time.Duration
}

type PostgresConfig struct {
	Host            string
	User            string
	Password        string
	DbName          string
	SslMode         string
	Port            int
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host               string
	Port               int
	Password           string
	Db                 int
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	PoolSize           int
	PoolTimeout        int
	IdleCheckFrequency int
}
type OtpConfig struct {
	Digits     int
	ExpireTime time.Duration
	Limiter    time.Duration
}

type LoggerConfig struct {
	FilePath string
	Encoding string
	Level    string
	Logger   string
}

func findConfig(name string) (p string) {
	switch name {
	case "dev":
		return "../config/config-development.yml" // we use .. because it will find conf from cnd file
	case "dep":
		return "../config/config-production.yml"
	}
	return "../config/config-development.yml"
}

func loadConfig(filePath, fileType string) (v *viper.Viper, err error) {
	v = viper.New()
	v.SetConfigName(filePath)
	v.SetConfigType(fileType)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return v, nil
}

func parsConfig(v *viper.Viper) (cfg *Config, err error) {
	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func GetConfig() *Config {
	path := findConfig("dev")
	v, err := loadConfig(path, "yml")
	if err != nil {
		return nil
	}
	cfg, err := parsConfig(v)
	if err != nil {
		return nil
	}
	return cfg
}
