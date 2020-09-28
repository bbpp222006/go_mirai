package main

import "go_mirai/utils"

func main()  {
	output_flow:=make(chan string,100)
	go utils.Login(output_flow)
	for {
		println(<-output_flow)
	}
}