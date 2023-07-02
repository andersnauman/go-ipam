package config

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

type settings struct {
	Webserver webserver `json:"webserver"`
	SQL       sql       `json:"sql"`
}

type webserver struct {
	IP   net.IP `json:"ip"`
	Port int    `json:"port"`
}

type sql struct {
	User   string            `json:"user"`
	Passwd string            `json:"pass"`
	Net    string            `json:"net"`
	Addr   string            `json:"addr"`
	DBName string            `json:"dbname"`
	Params map[string]string `json:"params"`
}

var Settings settings

// init loads all the values into Settings
func init() {
	var file []byte
	var err error

	if _, err = os.Stat("/etc/settings.json"); err == nil {
		file, err = os.ReadFile("/etc/settings.json")
	} else {
		file, err = os.ReadFile("config/settings.json")
	}

	if err != nil {
		log.Printf("[!] " + err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(file, &Settings)
	if err != nil {
		log.Printf("[!] " + err.Error())
		os.Exit(1)
	}
}
