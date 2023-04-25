package model

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppRest            AppRest
	AppGrpc            AppGrpc
	SendingRetry       int
	SendingInterval    time.Duration
	ConnectingRetry    int
	ConnectingInterval time.Duration
}

type AppRest struct {
	Host     string
	Port     string
	Username string
	Password string
}

type AppGrpc struct {
	Host string
	Port string
	Key  string
}

const (
	EnvDefaultPath = "/etc/app/share/settings/"
	EnvFilename    = ".env"
)

func ReadConfig() (Config, error) {
	if err := godotenv.Load(fmt.Sprintf("%s%s", EnvDefaultPath, EnvFilename)); err != nil {
		if err := godotenv.Load(EnvFilename); err != nil {
			return Config{}, err
		}
	}

	sendingRetryCfg := os.Getenv("SENDING_RETRY")
	sendingRetry, err := strconv.Atoi(sendingRetryCfg)
	if err != nil {
		return Config{}, err
	}

	sendingIntervalCfg := os.Getenv("SENDING_INTERVAL")
	sendingInterval, err := strconv.Atoi(sendingIntervalCfg)
	if err != nil {
		return Config{}, err
	}

	connectingRetryCfg := os.Getenv("CONNECTING_RETRY")
	connectingRetry, err := strconv.Atoi(connectingRetryCfg)
	if err != nil {
		return Config{}, err
	}

	connectingIntervalCfg := os.Getenv("CONNECTING_INTERVAL")
	connectingInterval, err := strconv.Atoi(connectingIntervalCfg)
	if err != nil {
		return Config{}, err
	}

	return Config{
		AppRest: AppRest{
			Host:     os.Getenv("APP_REST_HOST"),
			Port:     os.Getenv("APP_REST_PORT"),
			Username: os.Getenv("APP_REST_USER"),
			Password: os.Getenv("APP_REST_PASSWORD"),
		},
		AppGrpc: AppGrpc{
			Host: os.Getenv("APP_GRPC_HOST"),
			Port: os.Getenv("APP_GRPC_PORT"),
			Key:  os.Getenv("APP_GRPC_KEY"),
		},
		SendingRetry:       sendingRetry,
		SendingInterval:    time.Second * time.Duration(sendingInterval),
		ConnectingRetry:    connectingRetry,
		ConnectingInterval: time.Second * time.Duration(connectingInterval),
	}, nil

}
