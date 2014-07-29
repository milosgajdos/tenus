package main

import (
	"fmt"
	"log"
	"net"

	"github.com/milosgajdos83/tenus"
)

func main() {
	// CREATE BRIDGE AND BRING IT UP
	br, err := tenus.NewBridgeWithName("vethbridge")
	if err != nil {
		log.Fatal(err)
	}

	brIp, brIpNet, err := net.ParseCIDR("10.0.41.1/16")
	if err != nil {
		log.Fatal(err)
	}

	if err := br.SetLinkIp(brIp, brIpNet); err != nil {
		fmt.Println(err)
	}

	if err = br.SetLinkUp(); err != nil {
		fmt.Println(err)
	}

	// CREATE VETH PAIR
	veth, err := tenus.NewVethPairWithOptions("myveth01", tenus.VethOptions{PeerName: "myveth02"})
	if err != nil {
		log.Fatal(err)
	}

	// ASSIGN IP ADDRESS TO THE HOST VETH INTERFACE
	vethHostIp, vethHostIpNet, err := net.ParseCIDR("10.0.41.2/16")
	if err != nil {
		log.Fatal(err)
	}

	if err := veth.SetLinkIp(vethHostIp, vethHostIpNet); err != nil {
		fmt.Println(err)
	}

	// ADD MYVETH01 INTERFACE TO THE MYBRIDGE BRIDGE
	myveth01, err := net.InterfaceByName("myveth01")
	if err != nil {
		log.Fatal(err)
	}

	if err = br.AddSlaveIfc(myveth01); err != nil {
		fmt.Println(err)
	}

	if err = veth.SetLinkUp(); err != nil {
		fmt.Println(err)
	}

	// PASS VETH PEER INTERFACE TO A RUNNING DOCKER BY PID
	pid, err := tenus.DockerPidByName("vethdckr", "/var/run/docker.sock")
	if err != nil {
		fmt.Println(err)
	}

	if err := veth.SetPeerLinkNsPid(pid); err != nil {
		log.Fatal(err)
	}

	// ALLOCATE AND SET IP FOR THE NEW DOCKER INTERFACE
	vethGuestIp, vethGuestIpNet, err := net.ParseCIDR("10.0.41.5/16")
	if err != nil {
		log.Fatal(err)
	}

	if err := veth.SetPeerLinkNetInNs(pid, vethGuestIp, vethGuestIpNet, nil); err != nil {
		log.Fatal(err)
	}
}
