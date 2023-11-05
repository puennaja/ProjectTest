package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppConfig appConfig
	Server    serverConfig
	Service   serviceConfig
	MongoDB   mongodbConfig
}

type appConfig struct {
	Service string `envconfig:"SERVICE"`
	Version string `envconfig:"APP_VERSION"`
}

type serverConfig struct {
	Host string `envconfig:"SERVER_HOST" default:"0.0.0.0"`
	Port string `envconfig:"SERVER_PORT" default:"8443"`
}

type serviceConfig struct {
	AuthSvc struct {
		AccessTokenExpire  time.Duration `envconfig:"AUTH_ACCESS_TOKEN_EXPIRE"`
		RefreshTokenExpire time.Duration `envconfig:"AUTH_REFRESH_TOKEN_EXPIRE"`
		Secret             string        `envconfig:"AUTH_SECRET"`
		ModelPath          string        `envconfig:"AUTH_MODEL_PATH"`
		PolicyPath         string        `envconfig:"AUTH_POLICY_PATH"`
	}
}

type mongodbConfig struct {
	Timeout  time.Duration `envconfig:"MONGODB_TIMEOUT"`
	URI      string        `envconfig:"MONGODB_URI"`
	Database string        `envconfig:"MONGODB_DB"`
}

var cfg Config

func Init() {
	_ = godotenv.Load()
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("read env error %v", err)
	}
}

func GetConfig() Config {
	return cfg
}

func (conf serverConfig) ServerURL() string {
	return "https://" + conf.Host + ":" + conf.Port
}
