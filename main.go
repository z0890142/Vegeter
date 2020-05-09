package main

import (
	"Vegeter/helper/Comman"
	"Vegeter/helper/DB"
	"Vegeter/router"
	"log"
	"net/http"
	"strings"

	"github.com/rs/cors"
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
	DB.CreateDbConn("mysql", "developer:smap01@tcp(34.80.172.124:3306)/vegeter", Log)
	router := router.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,

		AllowedHeaders: []string{"Authorization", "Content-Type", "Access-Control-Allow-Origin"},
		// Enable Debugging for testing, consider disabling in production
		AllowedMethods: []string{"GET", "UPDATE", "PUT", "POST", "DELETE"},
	})

	log.Fatal(http.ListenAndServe(":8088", c.Handler(router)))

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
