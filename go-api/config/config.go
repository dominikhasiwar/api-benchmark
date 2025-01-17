package config

import (
	"strconv"

	"github.com/Energie-Burgenland/ausaestung-info/utils/env"
)

type Config struct {
	Port               int
	AWSRegion          string
	TableName          string
	Endpoint           string
	AuthTenantId       string
	AuthClientId       string
	AuthScope          string
	AllowedCorsOrigins string
	SwaggerBasePath    string
}

func GetConfig() Config {
	return Config{
		Port:               parseEnvToInt("BE_PORT", "8080"),
		AWSRegion:          env.GetEnv("BE_AWS_REGION", "us-east-1"),
		TableName:          env.GetEnv("BE_AWS_TABLE_NAME", ""),
		Endpoint:           env.GetEnv("BE_AWS_ENDPOINT", ""),
		AuthTenantId:       env.GetEnv("BE_AUTH_TENANT_ID", ""),
		AuthClientId:       env.GetEnv("BE_AUTH_CLIENT_ID", ""),
		AuthScope:          env.GetEnv("BE_AUTH_SCOPE", ""),
		AllowedCorsOrigins: env.GetEnv("BE_CORS_ALLOWED_ORIGINS", ""),
		SwaggerBasePath:    env.GetEnv("BE_SWAGGER_BASE_PATH", "/"),
	}
}

func parseEnvToInt(key, defaultValue string) int {
	value := env.GetEnv(key, defaultValue)
	num, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return num
}
