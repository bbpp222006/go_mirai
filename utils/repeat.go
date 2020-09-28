package utils

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strconv"
)

func Repeat(input_str string)  {
	a:=input_str
	println("接收:"+a)
	message_chain := gjson.Get(a, "messageChain")
	sender_type:= gjson.Get(a, "type")
	var sender_qq int64
	if gjson.Get(a, "sender.group.id").Exists(){
		sender_qq = gjson.Get(a, "sender.group.id").Int()
	}else {
		sender_qq = gjson.Get(a, "sender.id").Int()
	}

	value, _ := sjson.Set("", "sessionKey", session_key)
	value, _ = sjson.Set(value, "target", sender_qq)
	for i, name := range message_chain.Array()[1:] {
		value, _ = sjson.SetRaw(value, "messageChain."+strconv.Itoa(i),name.String())
	}
	println("/send"+sender_type.String())
	println("发送:"+value)

	r:=Request_post(mirai_api_http_locate+"/send"+sender_type.String(),value)
	println("服务器返回:"+r)


}