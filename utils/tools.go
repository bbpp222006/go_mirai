package utils

import (
	"io/ioutil"
	"net/http"
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