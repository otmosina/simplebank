package util

import "github.com/spf13/viper"

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`      //= "postgres"
	DBSource      string `mapstructure:"DB_SOURCE"`      //= "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable"
	ServerAddress string `mapstructure:"SERVER_ADDRESS"` //= "0.0.0.0:8080"
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
