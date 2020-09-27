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

func load_ini(ini_path string) *ini.File  {
	cfg, err := ini.Load(ini_path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	return cfg
}

func get_session_key(host string,auth_key string) string  {
	auth_key_json, _ := sjson.Set("", "authKey", auth_key)
	a:=Request_post(host+"/auth",auth_key_json)
	session_key := gjson.Get(a, "session").Str
	return session_key
}

func verify_session_key(host string,session_key string,qq_num string,session_exist_flag chan bool)  {
	mapD := map[string]string{"sessionKey": session_key, "qq": qq_num}
	mapB, _ := json.Marshal(mapD)
	//verify_session_key_json, _ := sjson.Set("", "sessionKey", session_key)
	a:=Request_post(host+"/verify",string(mapB))
	if  gjson.Get(a, "code").Int()==0{
		session_exist_flag<- true
	}
}

func begin_ws_listen(host string,session_key string,session_exist_flag chan bool,output_flow chan string)  {
	<-session_exist_flag
	c, _, err := websocket.DefaultDialer.Dial("ws://"+host+"/all?sessionKey="+session_key, nil)
	if err != nil {
		println("dial:", err)
	}
	defer c.Close()

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
	session_exist_flag:=make(chan bool)
	cfg := load_ini("config.ini")
	//fmt.Println("login:", cfg.Section("login").Key("qq").String())
	host:= cfg.Section("login").Key("mirai_api_http_locate").String()
	authKey:= cfg.Section("login").Key("authKey").String()
	qq_num:= cfg.Section("login").Key("qq_num").String()

	session_key := get_session_key(host,authKey)
	go verify_session_key(host,session_key,qq_num,session_exist_flag)
	go begin_ws_listen(host,session_key,session_exist_flag,output_flow)


	//value, _ := sjson.Set("", "name.asd", "Anderson")
	////println(value)
	//values := gjson.Get(value, "name.asd")
	//println(values.Exists())
	//println(host,authKey,qq_num)
}