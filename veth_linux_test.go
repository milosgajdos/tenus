package tenus

import (
	"net"
	"testing"
	"time"
)

type vethTest struct {
	hostIfc     string
	vethOptions VethOptions
}

func Test_NewVethPair(t *testing.T) {
	veth, err := NewVethPair()
	if err != nil {
		t.Fatalf("NewVethPair() failed to run: %s", err)
	}

	vethIfcName := veth.NetInterface().Name
	vethPeerName := veth.PeerNetInterface().Name

	tl := &testLink{}
	if err := tl.prepTestLink(vethIfcName, ""); err != nil {
		t.Skipf("NewVethPair test requries external command: %v", err)
	}

	if _, err := net.InterfaceByName(vethIfcName); err != nil {
		t.Fatalf("Could not find %s on the host: %s", vethIfcName, err)
	}

	if _, err := net.InterfaceByName(vethPeerName); err != nil {
		t.Fatalf("Could not find %s on the host: %s", vethPeerName, err)
	}

	testRes, err := linkInfo(vethIfcName, "veth")
	if err != nil {
		tl.teardown()
		t.Fatalf("Failed to list %s operation mode: %s", vethIfcName, err)
	}

	if testRes.linkType != "veth" {
		tl.teardown()
		t.Fatalf("NewVethPair() failed: expected linktype veth, returned %s", testRes.linkType)
	}

	if err := tl.teardown(); err != nil {
		t.Fatalf("testLink.teardown failed: %v", err)
	} else {
		time.Sleep(10 * time.Millisecond)
	}
}

var vethOptionTests = []vethTest{
	{"vethHost01", VethOptions{"vethGuest01"}},
	{"vethHost02", VethOptions{"vethGuest02"}},
}

func Test_NewVethPairWithOptions(t *testing.T) {
	for _, tt := range vethOptionTests {
		tl := &testLink{}

		if err := tl.prepTestLink(tt.hostIfc, ""); err != nil {
			t.Skipf("NewVlanLink test requries external command: %v", err)
		}

		veth, err := NewVethPairWithOptions(tt.hostIfc, tt.vethOptions)
		if err != nil {
			t.Fatalf("NewVethPairWithOptions(%s, %v) failed to run: %s", tt.hostIfc, tt.vethOptions, err)
		}

		if _, err := net.InterfaceByName(tt.hostIfc); err != nil {
			t.Fatalf("Could not find %s on the host: %s", tt.hostIfc, err)
		}

		if _, err := net.InterfaceByName(tt.vethOptions.PeerName); err != nil {
			t.Fatalf("Could not find %s on the host: %s", tt.vethOptions.PeerName, err)
		}

		vethIfcName := veth.NetInterface().Name
		if vethIfcName != tt.hostIfc {
			tl.teardown()
			t.Fatalf("NewVethPairWithOptions(%s, %v) failed: expected host ifc %s, returned %s",
				tt.hostIfc, tt.vethOptions, tt.hostIfc, vethIfcName)
		}

		vethPeerName := veth.PeerNetInterface().Name
		if vethPeerName != tt.vethOptions.PeerName {
			tl.teardown()
			t.Fatalf("NewVethPairWithOptions(%s, %v) failed: expected peer ifc %s, returned %s",
				tt.hostIfc, tt.vethOptions, tt.vethOptions.PeerName, vethPeerName)
		}

		testRes, err := linkInfo(tt.hostIfc, "veth")
		if err != nil {
			tl.teardown()
			t.Fatalf("Failed to list %s operation mode: %s", tt.hostIfc, err)
		}

		if testRes.linkType != "veth" {
			tl.teardown()
			t.Fatalf("NewVethPairWithOptions(%s, %v) failed: expected linktype veth, returned %s",
				tt.hostIfc, tt.vethOptions, testRes.linkType)
		}

		if err := tl.teardown(); err != nil {
			t.Fatalf("testLink.teardown failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}
