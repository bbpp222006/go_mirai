package main

import (
	"github.com/tidwall/sjson"
)




func main() {
	//value, _ := sjson.Set("", "name.asd", "Anderson")
	////println(value)
	//values := gjson.Get(value, "name.asd")
	//println(values.Exists())
	value, _ := sjson.Set("", "name", "Tom")
	println(value)
}


