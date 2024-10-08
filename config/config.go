package config

import (
	"os"

	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/baltazar")
	viper.AddConfigPath("$HOME/.baltazar")
	viper.AddConfigPath(".")

	viper.SetDefault("model.api_key", "OPENAI_COMPATIBLE_API_KEY")
	viper.SetDefault("model.api_base_url", "https://api.mistral.ai/v1/")
	viper.SetDefault("model.id", "mistral-large-2407")
	viper.SetDefault("model.temperature", 0.3)
	viper.SetDefault("model.max_tokens", 2048)
	viper.SetDefault("model.system_prompt", "You are a helpful assistant")
	viper.SetDefault("api.port", 8080)

	err := viper.ReadInConfig()
	if err != nil {
		os.Create("config.yaml")
		viper.WriteConfig()
	}
}
