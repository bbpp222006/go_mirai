package main

import "go_mirai/utils"


//每个模块都是个工厂,给它接上管道后就开始自动工作.
func main()  {
	output_flow:=make(chan string,100)
	default_flow:=make(chan string,100)
	vip_flow:=make(chan string,100)
	go utils.Login(output_flow)

	go utils.Vip_filter(output_flow,vip_flow,default_flow)

	utils.Repeat(default_flow)


}