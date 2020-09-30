package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"time"
)





func get_session_key(session_flow chan bool)   {
	for {
		if verify_session_key(){
			println("sessionkey校验成功")
			session_flow<-true
			return
		}else {
			println("sessionkey校验失败,重新获取中")

			auth_key_json, _ := sjson.Set("", "authKey", auth_key)
			//println(auth_key_json)
			a:=Request_post(mirai_api_http_locate+"/auth",auth_key_json)
			//println("服务器返回:"+a)
			session_key = gjson.Get(a, "session").Str
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
		println(a,"5s后再次尝试获取")
		time.Sleep(5*time.Second)
		return false
	}
}

func begin_ws_listen(session_exist_flag chan bool,output_flow  chan <-string)  {
	<-session_exist_flag
	c, _, err := websocket.DefaultDialer.Dial("ws://"+mirai_api_http_locate+"/message?sessionKey="+session_key, nil)
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

func Login(output_flow chan <-string){
	session_key_done:=make(chan bool)

	Init_global_var()

	go get_session_key(session_key_done)


	go begin_ws_listen(session_key_done,output_flow)

	//value, _ := sjson.Set("", "name.asd", "Anderson")
	////println(value)
	//values := gjson.Get(value, "name.asd")
	//println(values.Exists())
	//println(host,authKey,qq_num)
}