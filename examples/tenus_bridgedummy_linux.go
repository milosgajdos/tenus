package main

import (
	"fmt"
	"log"

	"github.com/milosgajdos83/tenus"
)

func main() {
	// Create a new network bridge
	br, err := tenus.NewBridgeWithName("mybridge")
	if err != nil {
		log.Fatal(err)
	}

	// Bring the bridge up
	if err = br.SetLinkUp(); err != nil {
		fmt.Println(err)
	}

	// Create a dummy link
	dl, err := tenus.NewLink("mydummylink")
	if err != nil {
		log.Fatal(err)
	}

	// Add the dummy link into bridge
	if err = br.AddSlaveIfc(dl.NetInterface()); err != nil {
		log.Fatal(err)
	}

	// Bring the dummy link up
	if err = dl.SetLinkUp(); err != nil {
		fmt.Println(err)
	}
}
