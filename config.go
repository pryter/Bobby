package main

import (
	"Bobby/cmd"
	"bytes"
	_ "embed"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

//go:embed bobby-notes.yaml
var configFile []byte

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

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer(configFile))

	if err != nil {
		log.Panic().Err(err).Msg("Unable to find bobby-notes")
	}

	err = viper.Unmarshal(&Configs)

	if err != nil {
		log.Panic().Err(err).Msg("Unable to unpack config file.")
	}
}
