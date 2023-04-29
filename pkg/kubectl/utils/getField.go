package utils

import (
	"github.com/spf13/viper"
)

/*
reference to :
https://github.com/spf13/viper
*/
func GetField(dirname string, filenameWithoutExtention string, field string) (string, error) {
	configs := viper.New()
	configs.SetConfigName(filenameWithoutExtention)
	configs.SetConfigType("yaml")
	configs.AddConfigPath(dirname)

	if err := configs.ReadInConfig(); err != nil {
		panic(err)
	}

	return configs.GetString(field), nil

}
