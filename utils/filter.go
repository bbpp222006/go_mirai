package utils

import "github.com/tidwall/gjson"

func Vip_filter(output_flow <-chan string,vip_flow chan <-string,default_flow  chan <-string)  {
	var sender_qq string
	for {
		origin_str:=<-output_flow
		if gjson.Get(origin_str, "sender.group.id").Exists() {
			sender_qq = gjson.Get(origin_str, "sender.group.id").String()
		} else {
			sender_qq = gjson.Get(origin_str, "sender.id").String()
		}
		if cfg.Section("auth").HasKey(sender_qq){
			vip_flow<-origin_str
		}else {
			default_flow<-origin_str
		}

	}
}
