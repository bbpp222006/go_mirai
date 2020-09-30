package utils

import (
	"github.com/imroc/req"
	"github.com/leekchan/accounting"
	"github.com/tidwall/gjson"
	"regexp"
	"strconv"
)

func Get_jita(repeat_flow <-chan string) {
	ac := accounting.Accounting{Symbol: "", Precision: 0}
	for {
		origin_str := <-repeat_flow
		result := is_jita(origin_str)
		all_flag := is_all_flag(origin_str)
		if result != "0" {
			Quick_reply(origin_str, "开始查询"+result+"相关价格,请稍后...")
			header := req.Header{
				"Host":             "www.ceve-market.org",
				"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:73.0) Gecko/20100101 Firefox/73.0",
				"Accept":           "application/json, text/javascript, */*; q=0.01",
				"Accept-Language":  "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
				"Accept-Encoding":  "gzip, deflate, br",
				"X-Requested-With": "XMLHttpRequest",
				"Connection":       "keep-alive",
				"Referer":          "https://www.ceve-market.org/home/",
			}
			param := req.Param{
				"search":   result,
				"regionid": 0,
				"tq":       0,
			}
			r, _ := req.Get("https://www.ceve-market.org/query/", header, param)
			content := Unidecode(r.String())
			println(content)
			StatusCode := strconv.Itoa(r.Response().StatusCode)
			if StatusCode != "200" {
				Quick_reply(origin_str, "查询失败,代码:"+StatusCode+"  我的网络环境可能不是很好,待会再试试?")
			} else {
				//message_chain := gjson.Get(content, "0")
				returnstr := "名称：卖/买\n"
				var buy_whole_price, sell_whole_price int64
				for i, item_info := range gjson.Parse(content).Array() {
					match, _ := regexp.MatchString("涂装", gjson.Get(item_info.String(), "typename").String())
					if (match || gjson.Get(item_info.String(), "sell").Int() == 0) && !all_flag {
						continue
					} else {
						returnstr += gjson.Get(item_info.String(), "typename").String()
						returnstr += "  " + ac.FormatMoney(gjson.Get(item_info.String(), "sell").Int())
						returnstr += "/" + ac.FormatMoney(gjson.Get(item_info.String(), "buy").Int()) + "\n"
					}
					buy_whole_price += gjson.Get(item_info.String(), "buy").Int()
					sell_whole_price += gjson.Get(item_info.String(), "sell").Int()
					if i>10 {
						break
					}

				}
				returnstr += "统计：  " + ac.FormatMoney(sell_whole_price) + "/" + ac.FormatMoney(buy_whole_price)
				if sell_whole_price == 0 && buy_whole_price == 0 {
					returnstr = "未查询到相关物品在售卖,或只查询到涂装,请检查名称输入是否有误\n 若要查看所有结果 在命令后加all即可, 例如:jita 三钛合金 all"
				}
				Quick_reply(origin_str, returnstr)
			}

		}
	}
}

func is_jita(recent_str string) string {
	content := gjson.Get(recent_str, "messageChain.1.text").String()
	//match, err := regexp.MatchString("jita (.*?)", content)
	r, _ := regexp.Compile("jita (.*)( all)?")
	search_name := r.FindStringSubmatch(content)
	if search_name == nil || len(search_name[1]) == 0 {
		return "0"
	} else {
		return search_name[1]
	}
}

func is_all_flag(recent_str string) bool {
	content := gjson.Get(recent_str, "messageChain.1.text").String()
	//match, err := regexp.MatchString("jita (.*?)", content)
	match, _ := regexp.MatchString("all", content)
	return match
}
