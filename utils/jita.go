package utils

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
	"regexp"
)

func Get_jita(repeat_flow chan string)  {
	origin_str:=<-repeat_flow
	result:=is_jita(origin_str)
	if result!="0"{

		header := req.Header {
			"Host": "www.ceve-market.org",
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:73.0) Gecko/20100101 Firefox/73.0",
			"Accept": "application/json, text/javascript, */*; q=0.01",
			"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
			"Accept-Encoding": "gzip, deflate, br",
			"X-Requested-With": "XMLHttpRequest",
			"Connection": "keep-alive",
			"Referer": "https://www.ceve-market.org/home/",
		}

		data := req.Param{
			"search": name,
			"regionid": 0,
			"tq": 0,
		}
	}

}

func is_jita(recent_str string)  string{
	content := gjson.Get(recent_str, "messageChain.text").String()
	//match, err := regexp.MatchString("jita (.*?)", content)
	r, _ := regexp.Compile("jita (.*)")
	search_name:=r.FindStringSubmatch(content)
	if len(search_name[1])==0{
		return "0"
	}else {
		return search_name[1]
	}
}





