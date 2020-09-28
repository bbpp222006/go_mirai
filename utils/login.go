package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"os"
)

var (
	qq_num string
	auth_key string
	mirai_api_http_locate string
	session_key string
)

func load_ini(ini_path string) *ini.File  {
	cfg, err := ini.Load(ini_path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	return cfg
}

func get_session_key(session_flow chan bool)   {
	for {
		if verify_session_key(){
			println("sessionkey校验成功")
			session_flow<-true
			return
		}else {
			println("sessionkey校验失败,重新获取中")

			auth_key_json, _ := sjson.Set("", "auth_key", auth_key)
			a:=Request_post(mirai_api_http_locate+"/auth",auth_key_json)

			session_key := gjson.Get(a, "session").Str
			cfg := load_ini("config.ini")
			cfg.Section("login").Key("session_key").SetValue(session_key)
			cfg.SaveTo("config.ini")
			println("sessionkey设置为"+session_key+", 已保存至config.ini")
		}
	}


}

func verify_session_key() bool {
	mapD := map[string]string{"sessionKey": session_key, "qq": qq_num}
	mapB, _ := json.Marshal(mapD)
	//verify_session_key_json, _ := sjson.Set("", "sessionKey", session_key)
	a:=Request_post(mirai_api_http_locate+"/verify",string(mapB))
	if  gjson.Get(a, "code").Int()==0{
		return true
	}else {
		return false
	}
}

func begin_ws_listen(session_exist_flag chan bool,output_flow chan string)  {
	<-session_exist_flag
	c, _, err := websocket.DefaultDialer.Dial("ws://"+mirai_api_http_locate+"/all?sessionKey="+session_key, nil)
	if err != nil {
		println("dial:", err)
	}
	defer c.Close()
	fmt.Printf("开始在%s上监听\n",mirai_api_http_locate)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			println("read:", err)
		}
		//println("recv: %s", message)
		output_flow<-string(message)
	}


}

func Login(output_flow chan string){
	session_key_done:=make(chan bool)
	cfg := load_ini("config.ini")
	mirai_api_http_locate= cfg.Section("login").Key("mirai_api_http_locate").String()
	auth_key= cfg.Section("login").Key("auth_key").String()
	qq_num= cfg.Section("login").Key("qq_num").String()
	session_key= cfg.Section("login").Key("session_key").String()

	go get_session_key(session_key_done)


	go begin_ws_listen(session_key_done,output_flow)

	//value, _ := sjson.Set("", "name.asd", "Anderson")
	////println(value)
	//values := gjson.Get(value, "name.asd")
	//println(values.Exists())
	//println(host,authKey,qq_num)
}