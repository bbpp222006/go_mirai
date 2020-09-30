package utils

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func Spy(input_str_chan <-chan string) {
	var sender_qq string
	output_str_chan := make(chan string, 100)
	go Repeat_simple(output_str_chan)
	for {
		origin_str:=<-input_str_chan
		if gjson.Get(origin_str, "sender.group.id").Exists() {
			sender_qq = gjson.Get(origin_str, "sender.group.id").String()
		} else {
			sender_qq = gjson.Get(origin_str, "sender.id").String()
		}

		if cfg.Section("spy").HasKey(sender_qq){
			modified_str:=modify_str(origin_str)
			output_str_chan <- modified_str
		}else {
			continue
		}
	}
}

func modify_str(origin_str string)  string{
	var modified_str string
	sender_name:=gjson.Get(origin_str,"sender.memberName").String()
	group_name:=gjson.Get(origin_str,"sender.group.name").String()
	target:=cfg.Section("spy").Key("admin").String()
	modified_str, _ = sjson.Set(origin_str, "type", "GroupMessage")
	modified_str, _ = sjson.Delete(modified_str, "sender.group")
	modified_str, _ = sjson.Set(modified_str, "sender.id",target)

	modified_str, _ = sjson.Set(modified_str,"messageChain.-1", map[string]interface{}{"type": "Plain", "text": "\n"+group_name+"."+sender_name})
	return modified_str
}