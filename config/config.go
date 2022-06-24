package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type db struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	DbName string `json:"dbname"`
	Usr    string `json:"usr"`
	Pwd    string `json:"pwd"`
}

var (
	Env           string
	Db            db
	HttpPort      int
	TokenTime     int64 //seconds
	SendCoins     map[string]struct{}
	MaxQueryCount int
)

//Load load config file
func init() {
	file, err := os.Open("config/config.json")
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	type Config struct {
		Env           string   `json:"env"`
		HttpPort      int      `json:"http_port"`
		Db            db       `json:"db"`
		TokenTime     int64    `json:"token_time"`
		SendCoins     []string `json:"send_coin"`
		MaxQueryCount int      `json:"max_query_count"`
	}
	all := &Config{}
	if err = json.NewDecoder(file).Decode(all); err != nil {
		log.Panic(err)
	}
	Env = all.Env
	Db = all.Db
	HttpPort = all.HttpPort
	TokenTime = all.TokenTime
	SendCoins = make(map[string]struct{})
	for _, coin := range all.SendCoins {
		SendCoins[strings.ToUpper(coin)] = struct{}{}
	}
	MaxQueryCount = all.MaxQueryCount
}
