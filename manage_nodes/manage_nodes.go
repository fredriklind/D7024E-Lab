package main

import (
	//server "github.com/fredriklind/D7024E-Lab/manage_nodes/server"
	"./server"
)

func main() {
	go server.StartWebServer()
	go server.StartAPI()
	//server.Do()
	block := make(chan bool)
	<-block
}
