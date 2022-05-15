package gateway

import (
	"context"
	"fmt"
	"gmqtt/common"
	"gmqtt/mqtt"
	"math/rand"
	"net/http"
)

func SendMsgHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//获取消息内容
	msg := r.FormValue("msg")
	topic := r.FormValue("topic")
	if msg == "" || topic == "" {
		fmt.Fprintf(w, "msg or topic is error ")
		return
	}
	msgStruct := common.MsgStruct{Topic: topic, Content: msg}
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TRACEID, rand.Int())
	err := mqtt.GMergeWorker.HandPushMsg(ctx, &msgStruct)
	if err != nil {
		fmt.Fprintf(w, "push fail")
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	fmt.Fprintf(w, "success")
}
