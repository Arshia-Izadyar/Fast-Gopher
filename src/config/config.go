package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

var ProjectConfig *Config

type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Logger   LoggerConfig
	Otp      OtpConfig
	JWT      JWTConfig
	Cors     CorseConfig
	Key      KeyConfig
}

type KeyConfig struct {
	Len int `yaml:"len"`
}

type CorseConfig struct {
	AllowOrigins string
}
type ServerConfig struct {
	Port        int
	WorkerCount int
	RunMode     string
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

func LoadConfig() {
	path := findConfig("dev")
	v, err := loadConfig(path, "yml")
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := parsConfig(v)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
	ProjectConfig = cfg
}

func GetConfig() *Config {
	return ProjectConfig
}
