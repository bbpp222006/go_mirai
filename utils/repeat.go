package utils

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strconv"
)

func Repeat(input_str_chan <-chan string)  {
	threshold:=1
	output_str_chan:=make(chan string,100)
	go repeat(output_str_chan)
	for {
		i:=0
		recent_str:=""
		for i<threshold{
			new_str:=<-input_str_chan
			println(new_str)
			if is_repeat(new_str,recent_str){
				i++
			}else {
				recent_str=new_str
			}
		}
		output_str_chan<-recent_str
	}
}

func is_repeat(new_str string,recent_str string) bool {
	if len(recent_str) == 0{
		return false
	}

	var new_content ,old_content string
	for _, content := range  gjson.Get(new_str, "messageChain").Array()[1:] {
		//if gjson.Get(content.String(), "url").Exists(){
		content_without_url, _ := sjson.Delete(content.String(), "url",)

		new_content+=content_without_url
	}

	for _, content := range gjson.Get(recent_str, "messageChain").Array()[1:] {
		content_without_url, _ := sjson.Delete(content.String(), "url",)
		old_content+=content_without_url
	}

	//println(new_content,old_content)
	if new_content==old_content{
		return true
	}else {
		return false
	}
}

func repeat(input_str_chan <-chan string) {
	for {
		a := <-input_str_chan
		println("接收:" + a)
		message_chain := gjson.Get(a, "messageChain")
		sender_type := gjson.Get(a, "type")
		var sender_qq int64
		if gjson.Get(a, "sender.group.id").Exists() {
			sender_qq = gjson.Get(a, "sender.group.id").Int()
		} else {
			sender_qq = gjson.Get(a, "sender.id").Int()
		}

		value, _ := sjson.Set("", "sessionKey", session_key)
		value, _ = sjson.Set(value, "target", sender_qq)
		for i, name := range message_chain.Array()[1:] {
			value, _ = sjson.SetRaw(value, "messageChain."+strconv.Itoa(i), name.String())
		}
		println("/send" + sender_type.String())
		println("发送:" + value)

		r := Request_post(mirai_api_http_locate+"/send"+sender_type.String(), value)
		println("服务器返回:" + r)
	}

}
