package main

import (
	"Bobby/cmd"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var Configs Config

type ConcurrentPoolConfig struct {
	MaxConcurrentTasks int `mapstructure:"max_concurrent_tasks"`
}

type Config struct {
	AppVersion   string `mapstructure:"app_version"`
	HTTPServices struct {
		Webhook cmd.HTTPServiceConfig `mapstructure:"webhook"`
	} `mapstructure:"http_services"`
	ConcurrentPool ConcurrentPoolConfig `mapstructure:"concurrent_pool"`
}

func init() {

	viper.SetConfigName("bobby-notes")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		log.Panic().Err(err).Msg("Unable to find bobby-notes")
	}

	err = viper.Unmarshal(&Configs)

	if err != nil {
		log.Panic().Err(err).Msg("Unable to unpack config file.")
	}
}
