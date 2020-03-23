package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// DefaultViper .
func DefaultViper() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config/")
	v.AddConfigPath("../config/")
	v.AddConfigPath("../../config/")
	v.AddConfigPath("../../../config/")
	v.AddConfigPath("../../../../config/")
	v.AddConfigPath("../../../../../config/")
	err := v.ReadInConfig()
	if err != nil {
		panic(errors.WithMessage(err, "get viper fail"))
	}
	return v
}
