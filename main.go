package main

import (
	"Vegeter/helper/Comman"
	"Vegeter/helper/DB"
	"Vegeter/router"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/heroku/x/hmetrics/onload"
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
	DB.CreateDbConn("mysql", viper.GetString("DB.connectString"), Log)

	router := router.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,

		AllowedHeaders: []string{"Authorization", "Content-Type", "Access-Control-Allow-Origin"},
		// Enable Debugging for testing, consider disabling in production
		AllowedMethods: []string{"GET", "UPDATE", "PUT", "POST", "DELETE"},
	})
	port := os.Getenv("PORT")
	// port := "80"

	log.Fatal(http.ListenAndServe(":"+port, c.Handler(router)))

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
