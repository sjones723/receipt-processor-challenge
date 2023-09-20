package main

//import "fmt"
import( 

)


func main(){
	//fmt.Println("makefile working")
	//local storage
	receiptStore = make(map[string]Receipt)

	server := NewAPIServer(":3000")
	server.Run()
}