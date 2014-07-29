package main

import (
	"fmt"
	"log"
	"net"

	"github.com/milosgajdos83/tenus"
)

func main() {
	macVlanHost, err := tenus.NewMacVlanLinkWithOptions("eth1", tenus.MacVlanOptions{Mode: "bridge", MacVlanDev: "macvlanHostIfc"})
	if err != nil {
		log.Fatal(err)
	}

	macVlanHostIp, macVlanHostIpNet, err := net.ParseCIDR("10.0.41.2/16")
	if err != nil {
		log.Fatal(err)
	}

	if err := macVlanHost.SetLinkIp(macVlanHostIp, macVlanHostIpNet); err != nil {
		fmt.Println(err)
	}

	if err = macVlanHost.SetLinkUp(); err != nil {
		fmt.Println(err)
	}

	macVlanDocker, err := tenus.NewMacVlanLinkWithOptions("eth1", tenus.MacVlanOptions{Mode: "bridge", MacVlanDev: "macvlanDckrIfc"})
	if err != nil {
		log.Fatal(err)
	}

	pid, err := tenus.DockerPidByName("mcvlandckr", "/var/run/docker.sock")
	if err != nil {
		log.Fatal(err)
	}

	if err := macVlanDocker.SetLinkNetNsPid(pid); err != nil {
		log.Fatal(err)
	}

	macVlanDckrIp, macVlanDckrIpNet, err := net.ParseCIDR("10.0.41.3/16")
	if err != nil {
		log.Fatal(err)
	}

	if err := macVlanDocker.SetLinkNetInNs(pid, macVlanDckrIp, macVlanDckrIpNet, nil); err != nil {
		log.Fatal(err)
	}
}
