package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"flag"

	"gopkg.in/yaml.v2"
)

func main() {
	configFile := "ezstub.yaml"
	flag.StringVar(&configFile, "c", configFile, "Configuration file")
	flag.Parse()

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
