package tenus

import (
	"bytes"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

type ifcNameTest struct {
	ifcName  string
	expected bool
}

var ifcNameTests = []ifcNameTest{
	{"link1", true},
	{"", false},
	{"a", false},
	{"abcdefghijklmnopqr", false},
	{"link\uF021", false},
	{"eth0.123", true},
}

func Test_NetInterfaceNameValid(t *testing.T) {
	for _, tt := range ifcNameTests {
		ret, _ := NetInterfaceNameValid(tt.ifcName)
		if ret != tt.expected {
			t.Errorf("NetInterfaceNameValid(%s): expected %v, returned %v", tt.ifcName, tt.expected, ret)
		}
	}
}

type ifcMacTest struct {
	testLink
	opts     LinkOptions
	testVal  string
	expected *net.Interface
}

// correct MAC Address will always parse into HardwareAddr
var hw, _ = net.ParseMAC("22:ce:e0:99:63:6f")
var ifcMacTests = []ifcMacTest{
	{testLink{name: "ifc01", linkType: "dummy"},
		LinkOptions{MacAddr: "22:ce:e0:99:63:6f"}, "22:ce:e0:99:63:6f",
		&net.Interface{Name: "ifc01", HardwareAddr: hw}},
	{testLink{name: "ifc02", linkType: "dummy"},
		LinkOptions{MacAddr: "26:2e:71:98:60:8f"}, "",
		nil},
	{testLink{name: "ifc03", linkType: "dummy"},
		LinkOptions{MacAddr: "fa:de:b0:99:52:1c"}, "randomstring",
		nil},
}

func Test_FindInterfaceByMacAddress(t *testing.T) {
	for _, tt := range ifcMacTests {
		tl := &testLink{}

		if err := tl.prepTestLink(tt.name, tt.linkType); err != nil {
			t.Skipf("InterfaceByMacAddress test requries external command: %v", err)
		}

		if err := tl.prepLinkOpts(LinkOptions{MacAddr: tt.opts.MacAddr}); err != nil {
			t.Skipf("InterfaceByMacAddress test requries external command: %v", err)
		}

		if err := tl.create(); err != nil {
			t.Fatalf("testLink.create failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}

		if err := tl.setup(); err != nil {
			t.Fatalf("testLink.setup failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}

		ifc, err := FindInterfaceByMacAddress(tt.testVal)
		if ifc != nil {
			if tt.expected != nil {
				if ifc.Name != tt.expected.Name || !bytes.Equal(ifc.HardwareAddr, tt.expected.HardwareAddr) {
					tl.teardown()
					t.Fatalf("FindInterfaceByMacAddress(%s): expected %v, returned %v",
						tt.testVal, tt.expected, ifc)
				}
			}

			if tt.expected == nil {
				tl.teardown()
				t.Fatalf("FindInterfaceByMacAddress(%s): expected %v, returned %v ",
					tt.testVal, tt.expected, ifc)
			}
		}

		if ifc == nil {
			if tt.expected != nil {
				tl.teardown()
				t.Fatalf("FindInterfaceByMacAddress(%s): expected %v, returned %v, error: %s",
					tt.testVal, tt.expected, ifc, err)
			}
		}

		if err := tl.teardown(); err != nil {
			t.Fatalf("testLink.teardown failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

type dockerPidTest struct {
	testDocker
	host     string
	expected int
}

var dockerPidTests = []dockerPidTest{
	{testDocker{name: "topper1", command: "/usr/bin/top"}, "/var/run/docker.sock", 1234},
	{testDocker{name: "topper2", command: "/usr/bin/top"}, "somehost.com:9011", 0},
	{testDocker{name: "topper3", command: "/usr/bin/top"}, "", 0},
}

func Test_DockerPidByName(t *testing.T) {
	for _, tt := range dockerPidTests {
		td := &testDocker{}

		if err := td.prepTestDocker(tt.name, tt.command); err != nil {
			t.Skipf("DockerPidByName test requries external command: %v", err)
		}

		if err := td.create(); err != nil {
			t.Fatalf("prepTestDocker.create failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}

		retPid, err := DockerPidByName(tt.name, tt.host)
		if retPid != 0 {
			if tt.expected != 0 {
				dockerPath, err := exec.LookPath("docker")
				if err != nil {
					t.Errorf("Unable to find docker in PATH: %s", err)
				}

				out, err := exec.Command(dockerPath, "inspect", "-f", "{{.State.Pid }}", tt.name).Output()
				if err != nil {
					td.teardown()
					t.Fatalf("Failed to run tt.testCmd.Output(): %s", err)
				}

				actualPid, err := strconv.Atoi(strings.TrimSpace(string(out)))
				if err != nil {
					td.teardown()
					t.Fatalf("Failed to run strconv.Atoi(strings.TrimSpace(string(%v))): %s",
						out, err)
				}

				if retPid != actualPid {
					td.teardown()
					t.Errorf("DockerPidByName(%s, %s): expected: %v, returned: %v",
						tt.name, tt.host, out, retPid)
				}
			}

			if tt.expected == 0 {
				td.teardown()
				t.Errorf("DockerPidByName(%s, %s): expected: %v, returned: %v",
					tt.name, tt.host, tt.expected, retPid)
			}
		}

		if retPid == 0 {
			if tt.expected != 0 {
				td.teardown()
				t.Errorf("DockerPidByName(%s, %s): expected: %v, returned: %v, error: %s",
					tt.name, tt.host, tt.expected, retPid, err)
			}
		}

		if err := td.teardown(); err != nil {
			t.Fatalf("testDocker.teardown failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

type netNsTest struct {
	pid      int
	testCmd  *exec.Cmd
	expected int
}

var netNsTests = []netNsTest{
	{1234, &exec.Cmd{
		Path: "",
		Args: []string{"docker", "inspect", "-f", "{{.State.Pid }}", "testdckr"}},
		1234},
	{0, &exec.Cmd{}, 0},
}

func Test_NetNsHandle(t *testing.T) {
	for _, tt := range netNsTests {
		td := &testDocker{}

		if err := td.prepTestDocker("testdckr", "/usr/bin/top"); err != nil {
			t.Skipf("NetNsHandle test requries external command: %v", err)
		}

		if err := td.create(); err != nil {
			t.Fatalf("prepTestDocker.create failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}

		if tt.pid != 0 {
			dockerPath, err := exec.LookPath("docker")
			if err != nil {
				t.Errorf("Unable to find docker in PATH: %s", err)
			}

			out, err := exec.Command(dockerPath, "inspect", "-f", "{{.State.Pid }}", "testdckr").Output()
			if err != nil {
				td.teardown()
				t.Fatalf("Failed to run tt.testCmd.Output(): %s", err)
			}

			actualPid, err := strconv.Atoi(strings.TrimSpace(string(out)))
			if err != nil {
				td.teardown()
				t.Fatalf("Failed to run strconv.Atoi(strings.TrimSpace(string(%v))): %s",
					out, err)
			}

			tt.pid = actualPid
		}

		nsFd, err := NetNsHandle(tt.pid)
		if nsFd == 0 {
			if tt.expected != 0 {
				td.teardown()
				t.Fatalf("NetNsHandle(%d): expected: non-zero, returned: %v, error: %s",
					tt.pid, nsFd, err)
			}
		}

		if nsFd != 0 {
			if tt.expected == 0 {
				td.teardown()
				t.Fatalf("NetNsHandle(%d): expected: %d, returned: %v",
					tt.pid, tt.expected, nsFd)
			}
		}

		if err := td.teardown(); err != nil {
			t.Fatalf("testDocker.teardown failed: %v", err)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
}
