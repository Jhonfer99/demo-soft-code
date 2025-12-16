package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	App struct {
		Name        string `mapstructure:"name"`
		Port        int    `mapstructure:"port"`
		Environment string `mapstructure:"environment"`
	} `mapstructure:"app"`
	
	Database struct {
		Port     int    `mapstructure:"port"`
		Host     string `mapstructure:"host"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
	} `mapstructure:"database"`
	
	Namespace string `mapstructure:"namespace"`
	Owner     string `mapstructure:"owner"`
}

func Load() *AppConfig {
	viper.SetConfigFile(".env")
	viper.ReadInConfig() 
	
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error loading config.yaml: %v", err)
	}
	
	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}
	
	return &cfg
}