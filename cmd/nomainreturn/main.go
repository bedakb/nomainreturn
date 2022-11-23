package main

import (
	"fmt"
	"log"

	"github.com/bedakb/nomainreturn/nomainreturn"
	"github.com/spf13/viper"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	singlechecker.Main(nomainreturn.NewAnalyzer(cfg))
}

func loadConfig() (nomainreturn.NoMainReturnConfig, error) {
	viper.SetConfigName(".nomainreturn")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.nomainreturn")
	viper.AddConfigPath(".")
	viper.SetDefault("allowPackages", nomainreturn.DefaultAllowPackages)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nomainreturn.NoMainReturnConfig{}, fmt.Errorf("failed to parse config: %w", err)
		}
	}

	var cfg nomainreturn.NoMainReturnConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nomainreturn.NoMainReturnConfig{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return cfg, nil
}
