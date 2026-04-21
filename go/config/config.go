package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	QdrantURL            string `mapstructure:"QDRANT_URL"`
	POSTGRES_URL         string `mapstructure:"POSTGRES_URL"`
	Port                 string `mapstructure:"GO_PORT"`
	OllamaURL            string `mapstructure:"OLLAMA_URL"`
	OllamaIntentModel    string `mapstructure:"OLLAMA_INTENT_MODEL"`
	ChromeProfileDir     string `mapstructure:"CHROME_PROFILE_DIR"`
	ChromeDefaultProfile string `mapstructure:"CHROME_DEFAULT_PROFILE"`
	ChromePath           string `mapstructure:"CHROME_PATH"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)

	return config, err

}
