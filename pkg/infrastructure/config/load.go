package config

import (
	"embed"
	"os"
	"strings"
)

// LoadDotenv is a helper function to load the dotenv file into environment variables
func LoadDotenv(embedFS embed.FS, resourcePaths map[string]string) error {
	configFilePath := resourcePaths["configs"] + "/" + os.Getenv("APP_ENV") + ".env"
	envs, err := embedFS.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	if len(envs) == 0 {
		return nil
	}

	lines := strings.Split(string(envs), "\n")
	for _, line := range lines {
		if line != "" {
			splits := strings.SplitN(line, "=", 2)
			if os.Getenv(splits[0]) == "" {
				if err := os.Setenv(splits[0], splits[1]); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
