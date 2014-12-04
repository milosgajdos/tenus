package tenus

import (
	"net"
	"testing"
	"time"
)

type macvtpTest struct {
	masterDev string
}

var macvtpTests = []macvtpTest{
	{"master01"},
}

func Test_NewMacVtapLink(t *testing.T) {
	for _, tt := range macvtpTests {
		tl := &testLink{}

		if err := tl.prepTestLink(tt.masterDev, "dummy"); err != nil {
			t.Skipf("NewMacVtapLink test requries external command: %v", err)
		}

		if err := tl.create(); err != nil {
			t.Fatalf("testLink.create failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}

		macvtp, err := NewMacVtapLink(tt.masterDev)
		if err != nil {
			t.Fatalf("NewMacVtapLink(%s) failed to run: %s", tt.masterDev, err)
		}

		mvtpName := macvtp.NetInterface().Name
		if _, err := net.InterfaceByName(mvtpName); err != nil {
			tl.teardown()
			t.Fatalf("Could not find %s on the host: %s", mvtpName, err)
		}

		testRes, err := linkInfo(mvtpName, "macvtap")
		if err != nil {
			tl.teardown()
			t.Fatalf("Failed to list %s operation mode: %s", mvtpName, err)
		}

		if testRes.linkType != "macvtap" {
			tl.teardown()
			t.Fatalf("NewMacVtapLink(%s) failed: expected macvtap, returned %s",
				tt.masterDev, testRes.linkType)
		}

		if testRes.linkData != "bridge" {
			tl.teardown()
			t.Fatalf("NewMacVtapLink(%s) failed: expected bridge, returned %s",
				tt.masterDev, testRes.linkData)
		}

		if err := tl.teardown(); err != nil {
			t.Fatalf("testLink.teardown failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}
