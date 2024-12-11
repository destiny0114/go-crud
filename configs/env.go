package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

type Env struct {
	DBHost string `mapstructure:"DATABASE_URL"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&env)

	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}

	return &env
}
