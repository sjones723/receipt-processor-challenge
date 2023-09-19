package main

//import "fmt"
import( 

)


func main(){
	//fmt.Println("makefile working")
	receiptStore = make(map[string]Receipt)

	server := NewAPIServer(":3000")
	server.Run()
}