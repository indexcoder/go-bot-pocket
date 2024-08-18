package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	TelegramToken     string
	PocketConsumerKey string
	AuthServerURL     string
	TelegramBotUrl    string `mapstructure:"bot_url"`
	DBPath            string `mapstructure:"db_file"`

	Messages Messages
}

type Messages struct {
	Errors
	Responses
}

type Errors struct {
	Default      string `mapstructure:"default"`
	InvalidUrl   string `mapstructure:"invalidUrl"`
	Unauthorized string `mapstructure:"unauthorized"`
	UnableToSave string `mapstructure:"unableToSave"`
}

type Responses struct {
	Start             string `mapstructure:"start"`
	AlreadyAuthorized string `mapstructure:"alreadyAuthorized"`
	SavedSuccessfully string `mapstructure:"savedSuccessfully"`
	UnknownCommand    string `mapstructure:"unknownCommand"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")
	viper.SetConfigName("main")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return nil, err
	}

	var cnf Config

	if err := viper.Unmarshal(&cnf); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.responses", &cnf.Messages.Responses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.errors", &cnf.Messages.Errors); err != nil {
		return nil, err
	}

	if err := parseEnv(&cnf); err != nil {
		return nil, err
	}

	return &cnf, nil
}

func parseEnv(cnf *Config) error {

	os.Setenv("TOKEN", "7215266589:AAGsvQ2x_Da5KhUPpx2LWp0dTySrRgovCzU")
	os.Setenv("CONSUMER_KEY", "111953-4f0a374f95750a63cab5a5e")
	os.Setenv("AUTH_SERVER_URL", "http://localhost/")

	if err := viper.BindEnv("token"); err != nil {
		return err
	}

	if err := viper.BindEnv("consumer_key"); err != nil {
		return err
	}

	if err := viper.BindEnv("auth_server_url"); err != nil {
		return err
	}

	cnf.TelegramToken = viper.GetString("token")
	cnf.PocketConsumerKey = viper.GetString("consumer_key")
	cnf.AuthServerURL = viper.GetString("consumer_key")

	return nil

}
