package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	Driver   string
}

type APIConfig struct {
	APIPort string
}

type TokenConfig struct {
	IssuerName       string `json:"issuerName"`
	JwtSignatureKey  []byte `json:"JwtSignatureKey"`
	JwtSigningMethod *jwt.SigningMethodHMAC
	JwtExpiresTime   time.Duration
}

type Config struct {
	DBConfig
	APIConfig
	TokenConfig
}

func (c *Config) ConfigConfiguration() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("missing env file %v", err.Error())
	}
	c.DBConfig = DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Driver:   os.Getenv("DB_DRIVER"),
	}

	c.APIConfig = APIConfig{os.Getenv("API_PORT")}

	tokenExpire, _ := strconv.Atoi(os.Getenv("TOKEN_EXPIRE"))
	c.TokenConfig = TokenConfig{
		IssuerName:       os.Getenv("TOKEN_ISSUE"),
		JwtSignatureKey:  []byte(os.Getenv("TOKEN_SECRET")),
		JwtSigningMethod: jwt.SigningMethodHS256,
		JwtExpiresTime:   time.Duration(tokenExpire) * time.Minute,
	}

	switch {
	case c.Host == "":
		c.Host = "localhost"
	case c.Port == "":
		c.Port = "5432"
	case c.User == "":
		c.User = "postgres"
	case c.Password == "":
		c.Password = "postgres"
	case c.Driver == "":
		c.Driver = "postgres"
	case c.Name == "":
		c.Name = "postgres"
	case c.APIPort == "":
		c.APIPort = "8080"
	case c.IssuerName == "":
		c.IssuerName = "Rizkyyullah"
	case c.JwtExpiresTime <= 0:
		c.JwtExpiresTime = 5
	}

	return nil
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := cfg.ConfigConfiguration(); err != nil {
		return nil, err
	}
	return cfg, nil
}
