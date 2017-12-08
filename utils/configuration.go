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
func InitExecutionContext(configFileName string, context *AppContext) (err error) {
	var file *os.File
	if file, err = os.Open(configFileName); err != nil {
		return err
	}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(context); err != nil {
		return err
	}
	defer file.Close()
	return nil
}
