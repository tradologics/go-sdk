package go_sdk

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	Version = "v0.1.2"
)

func InitEnv() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Print(".env file not found")
	}
}

func getEnvString(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

type TestConfig struct {
	SandboxURL   string
	SandboxToken string
}

var TestCfg *TestConfig

func GetTestConfig() *TestConfig {
	if TestCfg != nil {
		return TestCfg
	}

	InitEnv()

	TestCfg = &TestConfig{
		SandboxURL:   getEnvString("TEST_SANDBOX_URL", ""),
		SandboxToken: getEnvString("TEST_SANDBOX_TOKEN", ""),
	}

	return TestCfg
}
