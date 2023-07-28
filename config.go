package main

import (
	"encoding/json"
	"log"
	"os"
)

type User struct {
	Token   string `json:"token"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Config struct {
	Users    []User `json:"users"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

func LoadConfig() Config {
	data, err := os.ReadFile("/etc/gnoreply/config.json")
	if err != nil {
		log.Fatalln("Error reading the JSON file:", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln("Error unmarshalling JSON:", err)
	}

	return config
}
