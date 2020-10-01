package utils

import (
	"github.com/imroc/req"
	"github.com/leekchan/accounting"
	"github.com/tidwall/gjson"
	"math/big"
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
			StatusCode,name_jsonlist := get_name_jsonlist(result)
			//r, _ := req.Get("https://www.ceve-market.org/query/", header, param)
			//content := Unidecode(r.String())
			//println(content)
			//StatusCode := strconv.Itoa(r.Response().StatusCode)
			if StatusCode != "200" {
				Quick_reply(origin_str, "查询失败,代码:"+StatusCode+"  我的网络环境可能不是很好,待会再试试?")
			} else {
				//message_chain := gjson.Get(content, "0")
				returnstr := "名称：卖/买\n"
				var buy_whole_price, sell_whole_price float64
				for i, item_info := range gjson.Parse(name_jsonlist).Array() {
					match, _ := regexp.MatchString("涂装", gjson.Get(item_info.String(), "typename").String())
					sell,buy:=get_single_price(item_info.String())

					if (match || sell == 0) && !all_flag {
						continue
					} else {
						returnstr += gjson.Get(item_info.String(), "typename").String()
						returnstr += "  " + ac.FormatMoneyBigFloat(big.NewFloat(sell))
						returnstr += "/" + ac.FormatMoneyBigFloat(big.NewFloat(buy)) + "\n"
					}
					buy_whole_price += buy
					sell_whole_price += sell
					if i>10 {
						break
					}
				}
				returnstr += "统计：  " + ac.FormatMoneyBigFloat(big.NewFloat(sell_whole_price)) + "/" + ac.FormatMoneyBigFloat(big.NewFloat(buy_whole_price))
				if sell_whole_price == 0 && buy_whole_price == 0 {
					returnstr = "未查询到相关物品在售卖,或只查询到涂装,请检查名称输入是否有误\n 若要查看所有结果 在命令后加all即可, 例如:jita 三钛合金 all"
				}
				Quick_reply(origin_str, returnstr)
			}

		}
	}
}

func get_name_jsonlist(name string) (string,string) {
	header := req.Header{
		"Host":       "www.ceve-market.org",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36 Edg/85.0.564.63",
		"Connection": "keep-alive",
		"Referer":    "https://www.ceve-market.org/home/",
	}
	param := req.Param{
		"name": name,
	}
	r, _ := req.Post("https://www.ceve-market.org/api/searchname", header, param)
	content := r.String()
	StatusCode := strconv.Itoa(r.Response().StatusCode)
	return StatusCode,content
}

func get_single_price(typeid string) (float64,float64) {
	//ac := accounting.Accounting{Symbol: "", Precision: 1}
	header := req.Header{
		"Host":       "www.ceve-market.org",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.121 Safari/537.36 Edg/85.0.564.63",
		"Connection": "keep-alive",
		"Referer":    "https://www.ceve-market.org/home/",
	}
	//name := gjson.Get(typeid,"typename").String()
	returnstr := gjson.Get(typeid,"typeid").String()
	url_ := "https://www.ceve-market.org/api/market/region/10000002/type/" + returnstr + ".json"
	r, _ := req.Get(url_, header)
	content := r.String()
	//buy := ac.FormatMoneyBigFloat(big.NewFloat(gjson.Get(content, "buy.max").Float()))
	//sell :=  ac.FormatMoneyBigFloat(big.NewFloat(gjson.Get(content, "sell.min").Float()))
	//name += ":  " + sell
	//name += "/" + buy + "\n"
	return gjson.Get(content, "sell.min").Float(),gjson.Get(content, "buy.max").Float()
}


func is_jita(recent_str string) string {
 	content := gjson.Get(recent_str, "messageChain.1.text").String()
	//match, err := regexp.MatchString("jita (.*?)", content)
	r, _ := regexp.Compile("[^\\s]*")
	search_name := r.FindAllStringSubmatch(content,-1)
	if len(search_name) == 1 ||search_name[0][0] != "jita" || len(search_name[1][0]) == 0 {
		return "0"
	} else {
		return search_name[1][0]
	}
}

func is_all_flag(recent_str string) bool {
	content := gjson.Get(recent_str, "messageChain.1.text").String()
	//match, err := regexp.MatchString("jita (.*?)", content)
	match, _ := regexp.MatchString("all", content)
	return match
}
