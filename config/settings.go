package config

import (
	"github.com/spf13/viper"
)

type Settings struct {
	Mqtt     Mqtt               `mapstructure:"mqtt"`
	Monitors map[string]Monitor `mapstructure:"monitors"`
	Ddcutil  Ddcutil            `mapstructure:"ddcutil"`
}

type Mqtt struct {
	Broker string `mapstructure:"broker"`
	Port   int    `mapstructure:"port"`
	Topic  string `mapstructure:"topic"`
}

type Monitor struct {
	Serial string            `mapstructure:"serial"`
	Inputs map[string]string `mapstructure:"inputs"`
}

type Ddcutil struct {
	Bin string `mapstructure:"bin"`
}

func ParseSettings() *Settings {
	settings := Settings{}
	v := viper.New()
	v.SetConfigType("json")
	v.AddConfigPath("./")
	v.SetConfigName("config.json")

	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = v.Unmarshal(&settings)
	if err != nil {
		panic(err)
	}

	return &settings
}
