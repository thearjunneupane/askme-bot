package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	Bot_Token   string
	Bot_Prefix  string
	Gpt_api_key string
	Gpt_api_url string
	config      *configStruct
)

type configStruct struct {
	Bot_Token   string `json:"Bot_Token"`
	Bot_Prefix  string `json:"Bot_Prefix"`
	Gpt_api_key string `json:"Gpt_api_key"`
	Gpt_api_url string `json:"Gpt_api_url"`
}

func ReadConfig() error {
	fmt.Println("Reading Config File...!")

	file, err := os.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Bot_Token = config.Bot_Token
	Bot_Prefix = config.Bot_Prefix
	Gpt_api_key = config.Gpt_api_key
	Gpt_api_url = config.Gpt_api_url

	return nil

}
