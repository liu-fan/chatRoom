package mqtt

import (
	"net/http"
	"strings"
)

var warnProFun *warnPro

type warnPro struct{}

func (warn *warnPro) DingRobot(msg string) {
	url := "https://oapi.dingtalk.com/robot/send?access_token=552a8aea20f3f505d56a110db292735c3d474cbe8a1903ccabd5cc2051a7436a"
	s1 := "{\"msgtype\": \"text\",\"text\": {\"content\":\""
	s2 := "\"}}"
	var build strings.Builder
	build.WriteString(s1)
	build.WriteString(msg)
	build.WriteString(s2)
	s3 := build.String()
	payload := strings.NewReader(s3)
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Content-Type", "application/json")
	http.DefaultClient.Do(req)
}
