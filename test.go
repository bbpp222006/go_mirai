package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"os"
)


var cfg_add *ini.File

func load_ini(ini_path string) *ini.File  {
	cfg, err := ini.Load(ini_path)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	return cfg
}

func main() {
	//value, _ := sjson.Set("", "name.asd", "Anderson")
	////println(value)
	//values := gjson.Get(value, "name.asd")
	//println(values.Exists())
	//value, _ := sjson.Set("", "name", "Tom")
	//println(value)
	cfg := load_ini("config.ini")
	cfg_add = cfg
	println(cfg_add.Section("login").Key("mirai_api_http_locate").String())
}


