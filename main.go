package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	viper.SetDefault("application.db.max_idle", 5)
	viper.SetDefault("application.db.max_conn", 10)

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetConfigFile(`config.yml`)
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
	}
}

func main() {

	log.SetFormatter(&log.JSONFormatter{})

	app := NewApplication()

	if len(os.Args) < 2 {
		fmt.Println("add server or consumer parameter.")
		return
	}

	arg := os.Args[1]

	switch arg {
	case "server":
		app.StartHTTPServer()
	case "consumer":
		app.StartConsumer()
	}
}
