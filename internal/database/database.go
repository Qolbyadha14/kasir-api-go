package database

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func Init() (*Config, error) {
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	name := viper.GetString("DB_NAME")

	if host == "" || port == "" || user == "" || name == "" {
		return nil, fmt.Errorf("database configuration missing")
	}

	return &Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     name,
	}, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name)
}
