package main

import (
	"fmt"
	"gmqtt/config"
	"gmqtt/gateway"
	"gmqtt/mqtt"
	"net/http"
	"os"
)

func main() {
	//初始配置文件
	err := config.InitConfig("./config/tsconfig.json")

	if err != nil {
		goto ERR
	}
	err = mqtt.InitStats()
	if err != nil {
		goto ERR
	}
	mqtt.InitClientMgr()
	mqtt.InitMerger()
	http.HandleFunc("/sendMsg", gateway.SendMsgHandler)
	http.ListenAndServe(":8000", nil)

ERR:
	fmt.Println(err)
	os.Exit(-1)
}
