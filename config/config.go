package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	Version = "0.2.0"
)

type TestConfig struct {
	SandboxURL   string
	SandboxToken string
}

var TestCfg *TestConfig

func GetTestConfig(envPath ...string) *TestConfig {
	if TestCfg != nil {
		return TestCfg
	}

	if len(envPath) > 0 {
		InitEnv(envPath[0])
	} else {
		InitEnv()
	}

	TestCfg = &TestConfig{
		SandboxURL:   getEnvString("TEST_SANDBOX_URL", ""),
		SandboxToken: getEnvString("TEST_SANDBOX_TOKEN", ""),
	}

	return TestCfg
}

func InitEnv(envPath ...string) {
	if len(envPath) > 0 {
		if err := godotenv.Load(envPath[0]); err != nil {
			log.Print(".env file not found")
		}
	} else {
		if err := godotenv.Load(); err != nil {
			log.Print(".env file not found")
		}
	}
}

func getEnvString(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
