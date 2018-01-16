package configuration

import (
	"github.com/spf13/viper"
)

var (
	// ErtebotConfigPath is path of config file
	ErtebotConfigPath = "config.yaml"

	// ErtebotConfig is config of project
	ErtebotConfig *viper.Viper
)

// LoadConfig loads Ertebot's config file from ErtebotConfigPath
func LoadConfig() error {
	ErtebotConfig = viper.New()
	ErtebotConfig.SetConfigFile(ErtebotConfigPath)
	if err := ErtebotConfig.ReadInConfig(); err != nil {
		return err
	}

	ErtebotConfig.SetDefault("debug", true)

	return nil
}
