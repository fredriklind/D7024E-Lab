package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// start webserver and API

	// testing executing shell commands
	do()
}

// this file should contain the logic for the API-requests:
// startNewNode, nodeLeavesRing?, updateData (should return all nodes current stored data)?

// Starts a new Chord node in a Docker container. The node joins the ring.
func startNewNode() {
	cmd := exec.Command("docker", "images")
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print(string(out))
}

// execute shell commands from go code
func do() {

	//cmd := exec.Command("boot2docker", "start")
	cmd := exec.Command("docker", "images")
	out, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print(string(out))
}
