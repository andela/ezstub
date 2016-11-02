package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// configDir stores directory of the configuration file.
var configDir string

func main() {
	configFile := "ezstub.yaml"
	flag.StringVar(&configFile, "c", configFile, "Configuration file")
	flag.Parse()

	configDir = filepath.Dir(configFile)

	var config Config
	b, err := ioutil.ReadFile(configFile)
	exitIfErr(err)

	exitIfErr(yaml.Unmarshal(b, &config))

	server, err := NewServer(config)
	exitIfErr(err)

	log.Fatal(server.Start())
}

func exitIfErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
