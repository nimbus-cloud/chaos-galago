package main

import (
	"flag"
	"fmt"

	conf "chaos-galago/broker/config"
	utils "chaos-galago/broker/utils"
	webs "chaos-galago/broker/web_server"
)

// Options struct
type Options struct {
	ConfigPath string
}

var options Options

func init() {
	defaultConfigPath := utils.GetPath([]string{"assets", "config.json"})
	flag.StringVar(&options.ConfigPath, "c", defaultConfigPath, "use '-c' option to specify the config file path")

	flag.Parse()
}

func main() {
	var err error
	_, err = conf.LoadConfig(options.ConfigPath)
	if err != nil {
		panic(fmt.Sprintf("Error loading config file [%s]...", err.Error()))
	}

	server, err := webs.CreateServer()
	if err != nil {
		panic(fmt.Sprintf("Error creating server [%s]...", err.Error()))
	}

	server.Start()
}
