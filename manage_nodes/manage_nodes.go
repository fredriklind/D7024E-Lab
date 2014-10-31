package main

import (
	//server "github.com/fredriklind/D7024E-Lab/manage_nodes/server"
	//"./server"
	"fmt"
	"github.com/fsouza/go-dockerclient"
)

func main() {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	imgs, _ := client.ListImages(true)
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentId)
	}
	//go server.StartWebServer()
	//go server.StartAPI()
	//server.Do()
	//block := make(chan bool)
	//<-block
}
