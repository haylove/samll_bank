//  *@createTime    2022/3/20 19:23
//  *@author        hay&object
//  *@version       v1.0.0

package util

import (
	"github.com/spf13/viper"
	"time"
)

// Config store all configuration of application.
// The value are read by viper from a config file or environment variables.
type Config struct {
	DBDrive             string        `mapstructure:"DB_DRIVE"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

//LoadConfig load a configuration from a  file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	return
}
