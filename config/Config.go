package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	MQTTClientID         string `json:"MQTTClientID"`
	MaxMergerBatchSize   int    `json:"MaxMergerBatchSize"`
	MaxMergerDelay       int    `json:"MaxMergerDelay"`
	MergerChannelSize    int    `json:"MergerChannelSize"`
	MergerWorkerCount    int    `json:"MergerWorkerCount"`
	MergerPendingWarnNum int64  `json:"MergerPendingWarnNum"`
}

var GConfig *Config

func InitConfig(filePath string) (err error) {
	var (
		content []byte
		conf    Config
	)
	if content, err = ioutil.ReadFile(filePath); err != nil {
		return
	}

	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	GConfig = &conf
	return
}
