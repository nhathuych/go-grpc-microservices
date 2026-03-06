package config

import "github.com/nhathuych/api-gateway/pkg/util"

type Config struct {
	Port              string
	OrderServiceURL   string
	AccountServiceURL string
}

func LoadConfig() *Config {
	return &Config{
		Port:              util.GetEnv("PORT", "8080"),
		AccountServiceURL: util.GetEnv("ACCOUNT_SERVICE_URL", "http://localhost:8081"),
		OrderServiceURL:   util.GetEnv("ORDER_SERVICE_URL", "http://localhost:8082"),
	}
}
