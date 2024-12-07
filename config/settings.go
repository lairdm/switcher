package config

import (
	"github.com/spf13/viper"
)

type Settings struct {
	Mqtt     Mqtt               `mapstructure:"mqtt"`
	Monitors map[string]Monitor `mapstructure:"monitors"`
	Ddcutil  Ddcutil            `mapstructure:"ddcutil"`
	Xset     Xset               `mapstructure:"xset"`
	Amixer   Amixer             `mapstructure:"amixer"`
	XLock    XLock              `mapstructure:"xlock"`
	Lirc     Lirc               `mapstructure:"lirc"`
}

type Mqtt struct {
	Broker   string `mapstructure:"broker"`
	Port     int    `mapstructure:"port"`
	Topic    string `mapstructure:"topic"`
	ClientId string `mapstructure:"client_id"`
}

type Monitor struct {
	Serial  string            `mapstructure:"serial"`
	Inputs  map[string]string `mapstructure:"inputs"`
	Display string            `mapstructure:"display"`
}

type Ddcutil struct {
	Bin string `mapstructure:"bin"`
}

type Xset struct {
	Bin string `mapstructure:"bin"`
}

type Amixer struct {
	Bin string `mapstructure:"bin"`
}

type XLock struct {
	Bin string `mapstructure:"bin"`
}

type Lirc struct {
	Bin    string `mapstructure:"bin"`
	Remote string `mapstructure:"remote"`
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
