package main

import (
	"fmt"
	"log"
	"net"

	"github.com/milosgajdos83/tenus"
)

func main() {
	// CREATE BRIDGE AND BRING IT UP
	br, err := tenus.NewBridgeWithName("mybridge")
	if err != nil {
		log.Fatal(err)
	}

	ip, ipNet, err := net.ParseCIDR("10.0.41.1/16")
	if err != nil {
		log.Fatal(err)
	}

	if err := br.SetLinkIp(ip, ipNet); err != nil {
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

	// ASSIGN AN IP TO MYVETH01
	ip, ipNet, err = net.ParseCIDR("10.0.41.2/16")
	if err != nil {
		log.Fatal(err)
	}

	if err = veth.SetLinkIp(ip, ipNet); err != nil {
		fmt.Println(err)
	}

	// ASSIGN AN IP TO MYVETH02
	ip, ipNet, err = net.ParseCIDR("10.0.41.3/16")
	if err != nil {
		log.Fatal(err)
	}

	if err := veth.SetPeerLinkIp(ip, ipNet); err != nil {
		fmt.Println(err)
	}

	// ADD MYVETH01 INTERFACE TO THE MYBRIDGE BRIDGE AND BRING IT UP
	// we could also simply do myveth01 := veth.NetInterface()
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

	// ADD MYVETH02 INTERFACE TO THE MYBRIDGE BRIDGE AND BRING IT UP
	// we could also simply do myveth01 := veth.NetInterface()
	myveth02, err := net.InterfaceByName("myveth02")
	if err != nil {
		log.Fatal(err)
	}

	if err = br.AddSlaveIfc(myveth02); err != nil {
		fmt.Println(err)
	}

	if err = veth.SetPeerLinkUp(); err != nil {
		fmt.Println(err)
	}

	// CREATE MACVLAN INTERFACE AND BRING IT UP
	macvlan, err := tenus.NewMacVlanLinkWithOptions("eth0", tenus.MacVlanOptions{Mode: "bridge", MacVlanDev: "macvlan01"})
	if err != nil {
		log.Fatal(err)
	}

	if err := macvlan.SetLinkUp(); err != nil {
		fmt.Println(err)
	}

	// CREATE VLAN INTERFACE AND BRING IT UP
	vlan, err := tenus.NewVlanLinkWithOptions("eth1", tenus.VlanOptions{Id: 10, VlanDev: "vlan01"})
	if err != nil {
		log.Fatal(err)
	}

	if err = vlan.SetLinkUp(); err != nil {
		fmt.Println(err)
	}
}
