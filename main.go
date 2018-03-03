package main

import (
	"github.com/ChimeraCoder/anaconda"
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

const DEFAULT_CONFIG = "config.json"

func main(){
	config := loadConfig()
	api := anaconda.NewTwitterApiWithCredentials(
		config["access-token"].(string),
		config["your-access-token-secret"].(string),
		config["your-consumer-key"].(string),
		config["your-consumer-secret"].(string))
}

func loadConfig() map[string]interface{}{
	configFileName := getConfigFileName()
	checkConfigExist(configFileName)
	b, err := ioutil.ReadFile(configFileName)
	if err != nil {
		panic(err)
	}

	out := make(map[string]interface{})
	err = json.Unmarshal(b, &out)
	if err != nil{
		panic(err)
	}

	return out
}

func checkConfigExist(configName string) {
	if _, err := os.Stat(configName); os.IsNotExist(err) {
		panic(fmt.Sprintln("Couldn't load script properly, please check", configName, "filename"))
	}
}

func getConfigFileName() string{
	args := os.Args[1:]
	if len(args) > 0{
		return args[len(args) -1]
	}else{
		return DEFAULT_CONFIG
	}
}