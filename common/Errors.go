package common

import "errors"

var (
	ERR_MERGE_CHANNEL_FULL = errors.New("ERR_MERGE_CHANNEL_FULL")
	ERR_MQTT_MSG           = errors.New("ERR_MQTT_MSG")
)
