package main

import( 
//"fmt",
)

// OpenAPI annotations
// @title Receipt Processor
// @description A simple receipt processor
// @version 1.0.0
// @host localhost:3000
// @basePath /
// @schemes http
func main(){
	//fmt.Println("makefile working")
	//local storage
	receiptStore = make(map[string]Receipt)

	server := NewAPIServer(":3000")
	server.Run()
}