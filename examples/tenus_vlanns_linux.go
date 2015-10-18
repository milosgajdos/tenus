package main

import (
	"fmt"
	"log"
	"net"

	"github.com/milosgajdos83/tenus"
)

func main() {
	// CREATE VLAN HOST INTERFACE
	vlanDocker, err := tenus.NewVlanLinkWithOptions("eth1", tenus.VlanOptions{Dev: "vlanDckr", Id: 20})
	if err != nil {
		log.Fatal(err)
	}

	// PASS VLAN INTERFACE TO A RUNNING DOCKER BY PID
	pid, err := tenus.DockerPidByName("vlandckr", "/var/run/docker.sock")
	if err != nil {
		fmt.Println(err)
	}

	if err := vlanDocker.SetLinkNetNsPid(pid); err != nil {
		log.Fatal(err)
	}

	// ALLOCATE AND SET IP FOR THE NEW DOCKER INTERFACE
	vlanDckrIp, vlanDckrIpNet, err := net.ParseCIDR("10.1.41.3/16")
	if err != nil {
		log.Fatal(err)
	}

	if err := vlanDocker.SetLinkNetInNs(pid, vlanDckrIp, vlanDckrIpNet, nil); err != nil {
		log.Fatal(err)
	}
}
