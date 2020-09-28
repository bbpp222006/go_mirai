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
)


func load_ini(ini_path string) *ini.File  {
	cfg, err := ini.Load(ini_path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	return cfg
}

func init()  {
	cfg := load_ini("config.ini")
	mirai_api_http_locate= cfg.Section("login").Key("mirai_api_http_locate").String()
	auth_key= cfg.Section("login").Key("auth_key").String()
	qq_num= cfg.Section("login").Key("qq_num").String()
	session_key= cfg.Section("login").Key("session_key").String()
}
