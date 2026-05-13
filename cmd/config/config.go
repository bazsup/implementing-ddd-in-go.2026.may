package config

import "github.com/spf13/viper"

type Config struct {
	ServerPort        string
	WorkshopServerURL string
}

func LoadConfig() (Config, error) {
	viper.SetConfigFile(".env")
	viper.SetDefault("PORT", "8080")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, err
	}

	return Config{
		ServerPort:        ":" + viper.GetString("PORT"),
		WorkshopServerURL: viper.GetString("WORKSHOP_SERVER_URL"),
	}, nil
}
