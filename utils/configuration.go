package utils

import (
	"encoding/json"
	"os"
)

// AppContext is a configuration which is read from filw conf.json
type AppContext struct {
	APIURL      string `json:"APIURL"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	AccessToken string
	Expire      string
}

// InitExecutionContext is a method that reads settings from configuration file into special structure
func InitExecutionContext(context *AppContext) error {
	file, err := os.Open("myconf.json")
	if err != nil {
		return err
	}
	jsonParser := json.NewDecoder(file)
	err = jsonParser.Decode(context)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
