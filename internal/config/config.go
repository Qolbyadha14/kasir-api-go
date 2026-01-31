package config

import (
	"log"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
	URL  string `mapstructure:"url"`
}

type DatabaseConfig struct {
	URL          string `mapstructure:"url"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

var (
	cfg  *Config
	once sync.Once
)

func LoadConfig() *Config {
	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetConfigFile(".env")
	_ = v.ReadInConfig()

	// Set defaults from environment variables
	v.SetDefault("app.name", v.GetString("APP_NAME"))
	v.SetDefault("app.port", v.GetString("APP_PORT"))
	v.SetDefault("app.url", v.GetString("APP_URL"))
	v.SetDefault("database.url", v.GetString("DATABASE_URL"))
	v.SetDefault("database.max_open_conns", v.GetInt("DATABASE_MAX_OPEN_CONNS"))
	v.SetDefault("database.max_idle_conns", v.GetInt("DATABASE_MAX_IDLE_CONNS"))

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Println("Failed to unmarshal config:", err)
	}

	// Set fallback defaults if empty
	if config.App.Name == "" {
		config.App.Name = "kasir-api"
	}
	if config.App.Port == "" {
		config.App.Port = "8080"
	}
	if config.Database.MaxOpenConns == 0 {
		config.Database.MaxOpenConns = 25
	}
	if config.Database.MaxIdleConns == 0 {
		config.Database.MaxIdleConns = 25
	}

	return &config
}

func GetConfig() *Config {
	once.Do(func() {
		cfg = LoadConfig()
	})
	return cfg
}
