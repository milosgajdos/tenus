package tenus

import (
	"net"
	"testing"
	"time"
)

type macvlnTest struct {
	masterDev   string
	macvlanMode string
}

var macvlnTests = []macvlnTest{
	{"master01", "bridge"},
	{"master02", "private"},
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

		mvln, err := NewMacVlanLink(tt.masterDev, tt.macvlanMode)
		if err != nil {
			t.Fatalf("NewMacVlanLink(%s, %s) failed to run: %s", tt.masterDev, tt.macvlanMode, err)
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
			t.Fatalf("NewMacVlanLink(%s, %s) failed: expected macvlan, returned %s",
				tt.masterDev, tt.macvlanMode, testRes.linkType)
		}

		if testRes.linkData != tt.macvlanMode {
			tl.teardown()
			t.Fatalf("NewMacVlanLink(%s, %s) failed: expected %s, returned %s",
				tt.masterDev, tt.macvlanMode, tt.macvlanMode, testRes.linkData)
		}

		if err := tl.teardown(); err != nil {
			t.Fatalf("testLink.teardown failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}
