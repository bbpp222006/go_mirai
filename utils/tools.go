package utils

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func Request_post(url string,post_data string) string  {
	reader := strings.NewReader(post_data)
	req, _ := http.NewRequest("POST", "http://"+url, reader)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func Request_get(url string) string  {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://"+url, nil)
	if err != nil {
		// handle error
	}
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Add("Host", "tieba.baidu.com")
	//req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func Duplicate_chan(in chan string,out1 chan string,out2 chan string){
	message:=""
	for {
		message=<-in
		println(message)
		out1<-message
		out2<-message
	}
}

func Filter()  {
//	暂时还没想好这里该怎么实现一个通用的过滤器
}

func Quick_reply(origin_message string,str string)  {
	var sender_qq int64
	sender_type := gjson.Get(origin_message, "type")
	if gjson.Get(origin_message, "sender.group.id").Exists() {
		sender_qq = gjson.Get(origin_message, "sender.group.id").Int()
	} else {
		sender_qq = gjson.Get(origin_message, "sender.id").Int()
	}

	value, _ := sjson.Set("", "sessionKey", session_key)
	value, _ = sjson.Set(value, "target", sender_qq)
	value, _ = sjson.Set(value, "messageChain.0",  map[string]interface{}{"type":"Plain","text":str})


	println("发送:"+"/send" + sender_type.String() + value)
	r := Request_post(mirai_api_http_locate+"/send"+sender_type.String(), value)
	println("服务器返回:" + r)
}

func Unidecode(str string) string {

	name_str, _ := strconv.Unquote(strings.Replace(strconv.Quote(str), `\\u`, `\u`, -1))
	return name_str

}