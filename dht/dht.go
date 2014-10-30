package main

import (
	node "github.com/fredriklind/D7024E-Lab/dht/node"
	"os"
)

func main() {
	id := os.Args[1]
	ip := os.Args[2]
	transportPort := os.Args[3]
	nodeApiPort := os.Args[4]
	nodeDbPort := os.Args[5]
	node.NewLocalNode(&id, ip, transportPort, nodeApiPort, nodeDbPort)

	for {
		// do nothing
	}
}
