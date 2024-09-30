package config

import (
	"github.com/spf13/viper"
)

type Env struct {
    DBConnectionString string `mapstructure:"DB_CONNECTION_STRING"`
    SessionProvider string `mapstructure:"SESSION_MANAGER"`
    ServerPort string `mapstructure:"SERVER_PORT"`
    BrokerConnection string `mapstructure:"BROKER_CONNECTION_STRING"`
}

func NewEnv() *Env {
    env := Env{}

    viper.BindEnv("DB_CONNECTION_STRING")
    viper.SetDefault("DB_CONNECTION_STRING","mongodb://root:password@localhost:27017")
    viper.BindEnv("SESSION_PROVIDER")
    viper.SetDefault("SESSION_PROVIDER", "mongo")
    viper.BindEnv("SERVER_PORT")
    viper.SetDefault("SERVER_PORT", "4000")
    viper.BindEnv("BROKER_CONNECTION_STRING")
    viper.SetDefault("BROKER_CONNECTION_STRING", "amqp://guest:guest@broker.itsmelon.com:5672")

    env.DBConnectionString = viper.GetString("DB_CONNECTION_STRING")
    env.SessionProvider = viper.GetString("SESSION_PROVIDER")
    env.ServerPort = viper.GetString("SERVER_PORT")
    env.BrokerConnection = viper.GetString("BROKER_CONNECTION_STRING")


    return &env
}


