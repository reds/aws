package config

/*

 Load an aws config file. The file is json and has the following format:

{"service":{
	"host":"service.us-east-1.amazonaws.com",
	"secret":"secret key",
	"accessKeyId":"access id",
	"region":"us-east-1"
}}
*/

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config4 struct {
	Host        string `json:"host"`
	Region      string `json:"region"`
	Secret      string `json:"secret"`
	AccessKeyId string `json:"accesskeyid"`
}

func LoadConfig(fn, service string) (*Config4, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, errors.New("LoadConfig: " + err.Error())
	}
	fdata, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.New("LoadConfig: " + err.Error())
	}
	cfg := make(map[string]*Config4)
	err = json.Unmarshal(fdata, &cfg)
	if err != nil {
		return nil, errors.New("LoadConfig: " + err.Error())
	}
	return cfg[service], nil
}
