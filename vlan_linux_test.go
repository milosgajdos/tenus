package tenus

import (
	"net"
	"strconv"
	"testing"
	"time"
)

type vlnTest struct {
	masterDev string
	id        uint16
}

var vlnTests = []vlnTest{
	{"master01", 10},
	{"master02", 20},
}

func Test_NewVlanLink(t *testing.T) {
	for _, tt := range vlnTests {
		tl := &testLink{}

		if err := tl.prepTestLink(tt.masterDev, "dummy"); err != nil {
			t.Skipf("NewVlanLink test requries external command: %v", err)
		}

		if err := tl.create(); err != nil {
			t.Fatalf("testLink.create failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}

		vln, err := NewVlanLink(tt.masterDev, tt.id)
		if err != nil {
			t.Fatalf("NewVlanLink(%s, %s) failed to run: %s", tt.masterDev, tt.id, err)
		}

		vlnName := vln.NetInterface().Name
		if _, err := net.InterfaceByName(vlnName); err != nil {
			tl.teardown()
			t.Fatalf("Could not find %s on the host: %s", vlnName, err)
		}

		testRes, err := linkInfo(vlnName, "vlan")
		if err != nil {
			tl.teardown()
			t.Fatalf("Failed to list %s operation mode: %s", vlnName, err)
		}

		if testRes.linkType != "vlan" {
			tl.teardown()
			t.Fatalf("NewMacVlanLink(%s, %d) failed: expected vlan, returned %s",
				tt.masterDev, tt.id, testRes.linkType)
		}

		id, err := strconv.Atoi(testRes.linkData)
		if err != nil {
			tl.teardown()
			t.Fatalf("Failed to convert link data %s : %s", testRes.linkData, err)
		}

		if uint16(id) != tt.id {
			tl.teardown()
			t.Fatalf("NewMacVlanLink(%s, %d) failed: expected %d, returned %d",
				tt.masterDev, tt.id, tt.id, id)
		}

		if err := tl.teardown(); err != nil {
			t.Fatalf("testLink.teardown failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}
