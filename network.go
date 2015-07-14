package tenus

import (
	"github.com/opencontainers/runc/libcontainer/netlink"
)

type NetworkOptions struct {
	IpAddr string
	Gw     string
	Routes []netlink.Route
}
