package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"navaak/convertor/app"
	"os"
)

var configPath = flag.String("c", "config.json", "-c file/to/your/path.json")

func main() {
	config := loadConfig()
	a, _ := app.New(config)
	a.Run()
}

func loadConfig() app.Config {
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		return app.DefaultConfig
	}
	data, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	var config app.Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
