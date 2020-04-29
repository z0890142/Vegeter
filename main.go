package main

import (
	"Vegeter/helper/Comman"
	"Vegeter/router"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Log *logrus.Entry

func init() {
	Log, _ = Comman.LogInit("main", "Vegeter", logrus.DebugLevel)
	Log.Info("Vegeter version 0.0.0")
	InitConfig()

}

func main() {
	// DB.CreateDbConn("mysql", viper.GetString("DB.connectString"), Log)
	router := router.NewRouter()
	http.ListenAndServe(":80", router)
}

func InitConfig() {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()         // read in environment variables that match
	viper.SetEnvPrefix("gorush") // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("config")   // name of config file (without extension)
	viper.AddConfigPath("./config") // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err == nil {
		Log.Info("Using config file:", viper.ConfigFileUsed())
	}

}
