package main

//import "fmt"
//import "net/http"


func main(){
	//fmt.Println("makefile working")

	server := NewAPIServer(":3000")
	server.Run()
}