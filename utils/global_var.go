package utils

import (
	"fmt"
	"github.com/go-ini/ini"
	"os"
)

var (
	qq_num string
	auth_key string
	mirai_api_http_locate string
	session_key string
 	cfg *ini.File
	proxy_ip string
)


func load_ini(ini_path string) *ini.File  {
	cfg, err := ini.Load(ini_path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	return cfg
}

func Init_global_var()  {
	cfg = load_ini("config.ini")
	mirai_api_http_locate= cfg.Section("login").Key("mirai_api_http_locate").String()
	auth_key= cfg.Section("login").Key("auth_key").String()
	qq_num= cfg.Section("login").Key("qq_num").String()
	session_key= cfg.Section("login").Key("session_key").String()
	Change_proxy()
	return
}

func Change_proxy()  {
	//r, _ := req.Get("http://118.24.52.95/get/")
	//content := Unidecode(r.String())
	//proxy_ip=gjson.Get(content, "proxy").String()
	proxy_ip="34.92.243.136:9529"
	println("代理ip设置为"+proxy_ip)
}