package tenus

import (
	"net"
	"testing"
	"time"
)

type macvlnTest struct {
	masterDev string
}

var macvlnTests = []macvlnTest{
	{"master01"},
}

func Test_NewMacVlanLink(t *testing.T) {
	for _, tt := range macvlnTests {
		tl := &testLink{}

		if err := tl.prepTestLink(tt.masterDev, "dummy"); err != nil {
			t.Skipf("NewMacVlanLink test requries external command: %v", err)
		}

		if err := tl.create(); err != nil {
			t.Fatalf("testLink.create failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}

		mvln, err := NewMacVlanLink(tt.masterDev)
		if err != nil {
			t.Fatalf("NewMacVlanLink(%s) failed to run: %s", tt.masterDev, err)
		}

		mvlnName := mvln.NetInterface().Name
		if _, err := net.InterfaceByName(mvlnName); err != nil {
			tl.teardown()
			t.Fatalf("Could not find %s on the host: %s", mvlnName, err)
		}

		testRes, err := linkInfo(mvlnName, "macvlan")
		if err != nil {
			tl.teardown()
			t.Fatalf("Failed to list %s operation mode: %s", mvlnName, err)
		}

		if testRes.linkType != "macvlan" {
			tl.teardown()
			t.Fatalf("NewMacVlanLink(%s) failed: expected macvlan, returned %s",
				tt.masterDev, testRes.linkType)
		}

		if testRes.linkData != "bridge" {
			tl.teardown()
			t.Fatalf("NewMacVlanLink(%s) failed: expected bridge, returned %s",
				tt.masterDev, testRes.linkData)
		}

		if err := tl.teardown(); err != nil {
			t.Fatalf("testLink.teardown failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

type macvlnWithOptionsTest struct {
	masterDev string
	opts      *MacVlanOptions
}

var macvlnWithOptionsTests = []macvlnWithOptionsTest{
	{"master01", &MacVlanOptions{Dev: "test", MacAddr: "aa:aa:aa:aa:aa:aa", Mode: "bridge"}},
}

func Test_NewMacVlanLinkWithOptions(t *testing.T) {
	for _, tt := range macvlnWithOptionsTests {
		var iface *net.Interface

		tl := &testLink{}

		if err := tl.prepTestLink(tt.masterDev, "dummy"); err != nil {
			t.Skipf("NewMacVlanLinkWithOptions test requries external command: %v", err)
		}

		if err := tl.create(); err != nil {
			t.Fatalf("testLink.create failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}

		mvln, err := NewMacVlanLinkWithOptions(tt.masterDev, *tt.opts)
		if err != nil {
			t.Fatalf("NewMacVlanLinkWithOptions(%s, %s) failed to run: %s", tt.masterDev, *tt.opts, err)
		}

		iface = mvln.NetInterface()

		if iface.HardwareAddr.String() != tt.opts.MacAddr {
			tl.teardown()
			t.Fatalf("NewMacVlanLinkWithOptions(%s, %s) failed: expected %s, returned %s",
				tt.masterDev, *tt.opts, tt.opts.MacAddr, iface.HardwareAddr.String())
		}

		mvlnName := iface.Name
		if _, err := net.InterfaceByName(mvlnName); err != nil {
			tl.teardown()
			t.Fatalf("Could not find %s on the host: %s", mvlnName, err)
		}

		testRes, err := linkInfo(mvlnName, "macvlan")
		if err != nil {
			tl.teardown()
			t.Fatalf("Failed to list %s operation mode: %s", mvlnName, err)
		}

		if testRes.linkType != "macvlan" {
			tl.teardown()
			t.Fatalf("NewMacVlanLinkWithOptions(%s, %s) failed: expected macvlan, returned %s",
				tt.masterDev, *tt.opts, testRes.linkType)
		}

		if testRes.linkData != "bridge" {
			tl.teardown()
			t.Fatalf("NewMacVlanLinkWithOptions(%s, %s) failed: expected bridge, returned %s",
				tt.masterDev, *tt.opts, testRes.linkData)
		}

		if err := tl.teardown(); err != nil {
			t.Fatalf("testLink.teardown failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}
