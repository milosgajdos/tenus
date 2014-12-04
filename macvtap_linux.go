package tenus

import (
	"fmt"
	"net"

	"github.com/docker/libcontainer/netlink"
)

// MacVtaper embeds MacVlaner interface
type MacVtaper interface {
	MacVlaner
}

// MacVtapLink is MacVlanLink. It implements MacVtaper interface
type MacVtapLink struct {
	*MacVlanLink
}

// NewMacVtapLink creates macvtap network link
//
// It is equivalent of running:
//		ip link add name mvt${RANDOM STRING} link ${master interface} type macvtap mode ${mode}
// NewMacVtapLink returns MacVtaper which is initialized to a pointer of type MacVtapLink if the
// macvtap link was created successfully on the Linux host. Newly created link is assigned
// a random name starting with "mvt". It sets the macvlan mode to the parameter passed as argument.
// If incorrect network mode is passed as a paramter, it sets the macvlan mode to "bridge".
// It returns error if the link could not be created.
func NewMacVtapLink(masterDev string, mode string) (MacVtaper, error) {
	macVtapDev := makeNetInterfaceName("mvt")

	if ok, err := NetInterfaceNameValid(masterDev); !ok {
		return nil, err
	}

	if _, err := net.InterfaceByName(masterDev); err != nil {
		return nil, fmt.Errorf("Master MAC VTAP device %s does not exist on the host", masterDev)
	}

	if mode != "" {
		if _, ok := MacVlanModes[mode]; !ok {
			return nil, fmt.Errorf("Unsupported MacVtap mode specified: %s", mode)
		}
	} else {
		mode = "bridge"
	}

	if err := netlink.NetworkLinkAddMacVtap(masterDev, macVtapDev, mode); err != nil {
		return nil, err
	}

	macVtapIfc, err := net.InterfaceByName(macVtapDev)
	if err != nil {
		return nil, fmt.Errorf("Could not find the new interface: %s", err)
	}

	masterIfc, err := net.InterfaceByName(masterDev)
	if err != nil {
		return nil, fmt.Errorf("Could not find the new interface: %s", err)
	}

	return &MacVtapLink{
		MacVlanLink: &MacVlanLink{
			Link: Link{
				ifc: macVtapIfc,
			},
			masterIfc: masterIfc,
			mode:      mode,
		},
	}, nil
}

// NewMacVtapLinkWithOptions creates macvtap network link and can set some of its network parameters
// passed in as MacVlanOptions.
//
// It is equivalent of running:
// 		ip link add name ${macvlan name} link ${master interface} address ${macaddress} type macvtap mode ${mode}
// NewMacVtapLinkWithOptions returns MacVtaper which is initialized to a pointer of type MacVtapLink if the
// macvtap link was created successfully on the Linux host. It returns error if the macvtap link could not be created.
func NewMacVtapLinkWithOptions(masterDev string, opts MacVlanOptions) (MacVtaper, error) {
	macVtapDev := opts.MacVlanDev
	mode := opts.Mode
	macaddr := opts.MacAddr

	if ok, err := NetInterfaceNameValid(masterDev); !ok {
		return nil, err
	}

	if _, err := net.InterfaceByName(masterDev); err != nil {
		return nil, fmt.Errorf("Master MAC VLAN device %s does not exist on the host", masterDev)
	}

	if macVtapDev != "" {
		if ok, err := NetInterfaceNameValid(macVtapDev); !ok {
			return nil, err
		}

		if _, err := net.InterfaceByName(macVtapDev); err == nil {
			return nil, fmt.Errorf("MAC VATP device %s already assigned on the host", macVtapDev)
		}
	} else {
		macVtapDev = makeNetInterfaceName("mvt")
	}

	if mode != "" {
		if _, ok := MacVlanModes[mode]; !ok {
			return nil, fmt.Errorf("Unsupported MacVtap mode specified: %s", mode)
		}
	} else {
		mode = "bridge"
	}

	if err := netlink.NetworkLinkAddMacVtap(masterDev, macVtapDev, opts.Mode); err != nil {
		return nil, err
	}

	macVtapIfc, err := net.InterfaceByName(macVtapDev)
	if err != nil {
		return nil, fmt.Errorf("Could not find the new interface: %s", err)
	}

	if macaddr != "" {
		if _, err = net.ParseMAC(macaddr); err == nil {
			if err := netlink.NetworkSetMacAddress(macVtapIfc, macaddr); err != nil {
				if errDel := DeleteLink(macVtapIfc.Name); err != nil {
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

	return &MacVtapLink{
		MacVlanLink: &MacVlanLink{
			Link: Link{
				ifc: macVtapIfc,
			},
			masterIfc: masterIfc,
			mode:      mode,
		},
	}, nil
}
