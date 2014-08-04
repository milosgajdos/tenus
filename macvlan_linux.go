package tenus

import (
	"fmt"
	"net"

	"github.com/milosgajdos83/libcontainer-milosgajdos83/netlink"
)

// Supported macvlan modes by tenus package
var MacVlanModes = map[string]bool{
	"private": true,
	"vepa":    true,
	"bridge":  true,
}

// MacVlanOptions allows you to specify some options for macvlan link.
type MacVlanOptions struct {
	// macvlan device name
	MacVlanDev string
	// macvlan mode
	Mode string
	// MAC address
	MacAddr string
}

// MacVlaner embeds Linker interface and adds few more functions.
type MacVlaner interface {
	// Linker interface
	Linker
	// MasterNetInterface returns macvlan master network device
	MasterNetInterface() *net.Interface
	// Mode returns macvlan link's network mode
	Mode() string
}

// MacVlanLink is Link which has a master network device and operates in
// a given network mode. It implements MacVlaner interface.
type MacVlanLink struct {
	Link
	// Master device logical network interface
	masterIfc *net.Interface
	// macvlan operatio nmode
	mode string
}

// NewMacVlanLink creates macvlan network link
//
// It is equivalent of running:
//		ip link add name mc${RANDOM STRING} link ${master interface} type macvlan mode ${mode}
// NewMacVlanLink returns MacVlaner which is initialized to a pointer of type MacVlanLink if the
// macvlan link was created successfully on the Linux host. Newly created link is assigned
// a random name starting with "mc". It sets the macvlan mode the parameter passed as argument.
// If incorrect network mode is passed as a paramter, it sets the macvlan mode to "bridge".
// It returns error if the link could not be created.
func NewMacVlanLink(masterDev string, mode string) (MacVlaner, error) {
	macVlanDev := makeNetInterfaceName("mc")

	if ok, err := NetInterfaceNameValid(masterDev); !ok {
		return nil, err
	}

	if _, err := net.InterfaceByName(masterDev); err != nil {
		return nil, fmt.Errorf("Master MAC VLAN device %s does not exist on the host", masterDev)
	}

	if mode != "" {
		if _, ok := MacVlanModes[mode]; !ok {
			return nil, fmt.Errorf("Unsupported MacVlan mode specified: %s", mode)
		}
	} else {
		mode = "bridge"
	}

	if err := netlink.NetworkLinkAddMacVlan(masterDev, macVlanDev, mode); err != nil {
		return nil, err
	}

	macVlanIfc, err := net.InterfaceByName(macVlanDev)
	if err != nil {
		return nil, fmt.Errorf("Could not find the new interface: %s", err)
	}

	masterIfc, err := net.InterfaceByName(masterDev)
	if err != nil {
		return nil, fmt.Errorf("Could not find the new interface: %s", err)
	}

	return &MacVlanLink{
		Link: Link{
			ifc: macVlanIfc,
		},
		masterIfc: masterIfc,
		mode:      mode,
	}, nil
}

// NewMacVlanLinkWithOptions creates macvlan network link and sets som of its network parameters
// passed in as MacVlanOptions.
//
// It is equivalent of running:
// 		ip link add name ${macvlan name} link ${master interface} address ${macaddress} type macvlan mode ${mode}
// NewMacVlanLinkWithOptions returns MacVlaner which is initialized to a pointer of type MacVlanLink if the
// macvlan link was created successfully on the Linux host. It returns error if the macvlan link could not be created.
func NewMacVlanLinkWithOptions(masterDev string, opts MacVlanOptions) (MacVlaner, error) {
	macVlanDev := opts.MacVlanDev
	mode := opts.Mode
	macaddr := opts.MacAddr

	if ok, err := NetInterfaceNameValid(masterDev); !ok {
		return nil, err
	}

	if _, err := net.InterfaceByName(masterDev); err != nil {
		return nil, fmt.Errorf("Master MAC VLAN device %s does not exist on the host", masterDev)
	}

	if macVlanDev != "" {
		if ok, err := NetInterfaceNameValid(macVlanDev); !ok {
			return nil, err
		}

		if _, err := net.InterfaceByName(macVlanDev); err == nil {
			return nil, fmt.Errorf("MAC VLAN device %s already assigned on the host", macVlanDev)
		}
	} else {
		macVlanDev = makeNetInterfaceName("mc")
	}

	if mode != "" {
		if _, ok := MacVlanModes[mode]; !ok {
			return nil, fmt.Errorf("Unsupported MacVlan mode specified: %s", mode)
		}
	} else {
		mode = "bridge"
	}

	if err := netlink.NetworkLinkAddMacVlan(masterDev, macVlanDev, opts.Mode); err != nil {
		return nil, err
	}

	macVlanIfc, err := net.InterfaceByName(macVlanDev)
	if err != nil {
		return nil, fmt.Errorf("Could not find the new interface: %s", err)
	}

	if macaddr != "" {
		if _, err = net.ParseMAC(macaddr); err == nil {
			if err := netlink.NetworkSetMacAddress(macVlanIfc, macaddr); err != nil {
				if errDel := DeleteLink(macVlanIfc.Name); err != nil {
					return nil, fmt.Errorf("Incorrect options specified. Attempt to delete the link failed: %s",
						errDel)
				}
			}
		}
	}

	masterIfc, err := net.InterfaceByName(masterDev)
	if err != nil {
		return nil, fmt.Errorf("Could not find the new interface: %s", err)
	}

	return &MacVlanLink{
		Link: Link{
			ifc: macVlanIfc,
		},
		masterIfc: masterIfc,
		mode:      mode,
	}, nil
}

// NetInterface returns macvlan link's network interface
func (macvln *MacVlanLink) NetInterface() *net.Interface {
	return macvln.ifc
}

// MasterNetInterface returns macvlan link's master network interface
func (macvln *MacVlanLink) MasterNetInterface() *net.Interface {
	return macvln.masterIfc
}

// Mode returns macvlan link's network operation mode
func (macvln *MacVlanLink) Mode() string {
	return macvln.mode
}
