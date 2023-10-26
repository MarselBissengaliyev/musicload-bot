package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	TelegramApiToken string
	MaxDownloadTime  int64
	MaxVideoDuration int64
	BotUsername      string
	Debug            bool
}

func InitConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	err := viper.ReadInConfig()

	if err != nil {
		return nil, err
	}

	return &Config{
		TelegramApiToken: viper.GetString("TELEGRAM_APITOKEN"),
		MaxDownloadTime:  viper.GetInt64("MAX_DOWNLOAD_TIME"),
		MaxVideoDuration: viper.GetInt64("MAX_VIDEO_DURATION"),
		BotUsername:      viper.GetString("BOT_USERNAME"),
		Debug:            false,
	}, nil
}
