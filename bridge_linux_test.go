package tenus

import (
	"net"
	"testing"
	"time"
)

func Test_NewBridge(t *testing.T) {
	tl := &testLink{}

	br, err := NewBridge()
	if err != nil {
		t.Fatalf("NewBridge() failed to run: %s", err)
	}

	brName := br.NetInterface().Name
	if err := tl.prepTestLink(brName, "bridge"); err != nil {
		t.Skipf("NewBridge test requries external command: %v", err)
	}

	if _, err := net.InterfaceByName(brName); err != nil {
		tl.teardown()
		t.Fatalf("Could not find %s on the host: %s", brName, err)
	}

	testRes, err := linkInfo(brName, "bridge")
	if err != nil {
		tl.teardown()
		t.Fatalf("Failed to list %s operation mode: %s", brName, err)
	}

	if testRes.linkType != "bridge" {
		tl.teardown()
		t.Fatalf("NewBridge() failed: expected linktype bridge, returned %s", testRes.linkType)
	}

	if err := tl.teardown(); err != nil {
		t.Fatalf("testLink.teardown failed: %v", err)
	} else {
		time.Sleep(10 * time.Millisecond)
	}
}

func Test_NewBridgeWithName(t *testing.T) {
	brTests := []string{"br01", "br02", "br03"}
	for _, tt := range brTests {
		tl := &testLink{}

		_, err := NewBridgeWithName(tt)
		if err != nil {
			t.Fatalf("NewBridge(%s) failed to run: %s", tt, err)
		}

		if err := tl.prepTestLink(tt, "bridge"); err != nil {
			t.Skipf("test requries external command: %v", err)
		}

		if _, err := net.InterfaceByName(tt); err != nil {
			tl.teardown()
			t.Fatalf("Could not find %s on the host: %s", tt, err)
		}

		testRes, err := linkInfo(tt, "bridge")
		if err != nil {
			tl.teardown()
			t.Fatalf("Failed to list %s operation mode: %s", tt, err)
		}

		if testRes.linkType != "bridge" {
			tl.teardown()
			t.Fatalf("NewBridge() failed: expected linktype bridge, returned %s", testRes.linkType)
		}

		if err := tl.teardown(); err != nil {
			t.Fatalf("testIfcLink.teardown failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}
