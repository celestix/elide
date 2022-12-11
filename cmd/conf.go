package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

const elideConfigFile = "elide.yaml"

var (
	appId, protocol, port int
	apiHash, botToken     string
)

func loadConfig() {
	defer log.Println("[Elide][Config]: Loaded config vars...")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Elide: failed to read config:", err.Error())
		os.Exit(1)
	}
	appId = viper.GetInt("app-id")
	protocol = viper.GetInt("protocol")
	port = viper.GetInt("port")
	apiHash = viper.GetString("api-hash")
	botToken = viper.GetString("bot-token")
}

func createConfig() {
	viper.Set("app-id", 0)
	viper.Set("api-hash", "")
	viper.Set("bot-token", "")
	viper.Set("protocol", 0)
	viper.Set("port", 9093)
	if err := viper.WriteConfig(); err != nil {
		fmt.Println("Elide: failed to write config:", err.Error())
		os.Exit(1)
	}
}

func init() {
	viper.SetConfigFile(elideConfigFile)
}
