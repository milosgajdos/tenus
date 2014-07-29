package tenus

import (
	"net"
	"testing"
	"time"
)

func Test_NewLink(t *testing.T) {
	testLinks := []string{"ifc01", "ifc02", "ifc03"}
	for _, tt := range testLinks {
		tl := &testLink{}

		if err := tl.prepTestLink(tt, "dummy"); err != nil {
			t.Skipf("NewLink test requries external command: %v", err)
		}

		_, err := NewLink(tt)
		if err != nil {
			t.Fatalf("NewLink(%s) failed to run: %s", tt, err)
		}

		if _, err := net.InterfaceByName(tt); err != nil {
			tl.teardown()
			t.Fatalf("Could not find %s on the host: %s", tt, err)
		}

		if err := tl.teardown(); err != nil {
			t.Fatalf("testLink.teardown failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

type ifcLinkOptsTest struct {
	testLink
	opts     LinkOptions
	expected *net.Interface
}

var hwaddr, _ = net.ParseMAC("22:ce:e0:99:63:6f")
var ifcLinkOptsTests = []ifcLinkOptsTest{
	{testLink{name: "ifc01"},
		LinkOptions{MacAddr: "22:ce:e0:99:63:6f", MTU: 1400, Flags: net.FlagUp},
		&net.Interface{Name: "ifc01", MTU: 1400, Flags: net.FlagUp, HardwareAddr: hwaddr}},
	{testLink{name: "ifc02"},
		LinkOptions{},
		&net.Interface{Name: "ifc01"}},
	{testLink{name: "ifc03"},
		LinkOptions{MTU: -100},
		nil},
}

func Test_NewLinkWithOptions(t *testing.T) {
	for _, tt := range ifcLinkOptsTests {
		tl := &testLink{}

		if err := tl.prepTestLink(tt.name, "dummy"); err != nil {
			t.Skipf("NewLinkWithOptions test requries external command: %v", err)
		}

		_, err := NewLinkWithOptions(tt.name, tt.opts)
		if err != nil && tt.expected != nil {
			t.Fatalf("NewLinkWithOptions(%s, %v) failed to run: %s", tt.name, tt.opts, err)
		}

		ifc, err := net.InterfaceByName(tt.name)
		if err != nil && tt.expected != nil {
			tl.teardown()
			t.Fatalf("Could not find %s on the host: %s", tt.name, err)
		}

		if tt.expected != nil {
			if ifc.Name != tt.name {
				tl.teardown()
				t.Fatalf("NewLinkWithOptions(%s, %v) failed: expected %s, returned: %s",
					tt.name, tt.opts, tt.name, ifc.Name)
			}

			if tt.opts.MacAddr != "" {
				if ifc.HardwareAddr.String() != tt.opts.MacAddr {
					tl.teardown()
					t.Fatalf("NewLinkWithOptions(%s, %v) failed: expected %s, returned: %s",
						tt.name, tt.opts, tt.opts.MacAddr, ifc.HardwareAddr.String())
				}
			}

			if tt.opts.MTU != 0 {
				if ifc.MTU != tt.opts.MTU {
					tl.teardown()
					t.Fatalf("NewLinkWithOptions(%s, %v) failed: expected %d, returned: %d",
						tt.name, tt.opts, tt.opts.MTU, ifc.MTU)
				}
			}

			if tt.opts.Flags != 0 {
				if (ifc.Flags & tt.opts.Flags) != tt.opts.Flags {
					tl.teardown()
					t.Fatalf("NewLinkWithOptions(%s, %v) failed: expected %v, returned: %v",
						tt.name, tt.opts, tt.opts.Flags, ifc.Flags)
				}
			}

			if err := tl.teardown(); err != nil {
				t.Fatalf("testIfcLink.teardown failed: %v", err)
			} else {
				time.Sleep(10 * time.Millisecond)
			}
		}

		if tt.expected == nil && ifc != nil {
			tl.teardown()
			t.Fatalf("NewLinkWithOptions(%s, %v) failed. Expected: %v, Returned: %v",
				tt.name, tt.opts, tt.expected, ifc)
		}
	}
}

func Test_DeleteLink(t *testing.T) {
	testLinks := []string{"ifc01", "ifc02", "ifc03"}
	for _, tt := range testLinks {
		tl := &testLink{}

		if err := tl.prepTestLink(tt, "dummy"); err != nil {
			t.Skipf("DeleteLink test requries external command: %v", err)
		}

		if err := tl.create(); err != nil {
			t.Fatalf("testLink.setup failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}

		if err := DeleteLink(tt); err != nil {
			tl.teardown()
			t.Fatalf("Failed to delete %s interface: %s", tt, err)
		}

		i, _ := net.InterfaceByName(tt)
		if i != nil {
			tl.teardown()
			t.Fatalf("DeleteLink(%s) expected: nil, returned: %v", tt, i)
		}
	}
}
