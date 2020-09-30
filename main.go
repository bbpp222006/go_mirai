package main

import "go_mirai/utils"


//每个模块都是个工厂,给它接上管道后就开始自动工作.
func main()  {
	origin_flow:=make(chan string,100)
	default_flow:=make(chan string,100)
	vip_flow:=make(chan string,100)
	go utils.Login(origin_flow)

	main1:=make(chan string,100)
	spy_flow:=make(chan string,100)
	go utils.Duplicate_chan(origin_flow,main1,spy_flow)
	go utils.Spy(spy_flow)

	//
	//main2:=make(chan string,100)
	//spy_flow:=make(chan string,100)
	go utils.Vip_filter(main1,vip_flow,default_flow)

	jita_flow:=make(chan string,100)
	repeat_flow:=make(chan string,100)
	go utils.Duplicate_chan(default_flow,jita_flow,repeat_flow)

	go utils.Get_jita(jita_flow)
	go utils.Repeat(repeat_flow)

	for{
		<-vip_flow
	}
}