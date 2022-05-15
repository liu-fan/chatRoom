package mqtt

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"gmqtt/config"
	"log"
	"os"
	"time"
)

var GConnMgr *ClientMgr

type ClientMgr struct {
	MqttConn *mqtt.Client
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

//InitClientMgr 初始化MQTT连接
func InitClientMgr() {
	GConnMgr = &ClientMgr{}
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://broker.emqx.io:1883").SetClientID(config.GConfig.MQTTClientID)

	opts.SetKeepAlive(60 * time.Second)
	// 设置消息回调处理函数
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	GConnMgr.MqttConn = &c
	//// 订阅主题
	//if token := c.Subscribe("testtopic/#", 0, nil); token.Wait() && token.Error() != nil {
	//	fmt.Println(token.Error())
	//	os.Exit(1)
	//}

	// 发布消息
	//token := c.Publish("testtopic/1", 0, false, "Hello World")
	//token.Wait()

	//time.Sleep(6 * time.Second)
	//
	//// 取消订阅
	//if token := c.Unsubscribe("testtopic/#"); token.Wait() && token.Error() != nil {
	//	fmt.Println(token.Error())
	//	os.Exit(1)
	//}
	//
	//// 断开连接
	//c.Disconnect(250)
	//time.Sleep(1 * time.Second)
}
func (client *ClientMgr) PushMqtt(topic string, payload string) {
	c := *client.MqttConn
	c.Publish(topic, 0, false, payload)
}
