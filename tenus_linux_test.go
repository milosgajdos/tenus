package tenus

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type testEnv struct {
	createCmds   []*exec.Cmd
	setupCmds    []*exec.Cmd
	tearDownCmds []*exec.Cmd
}

func (te *testEnv) create() error {
	for _, cmd := range te.createCmds {
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (te *testEnv) setup() error {
	for _, cmd := range te.setupCmds {
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func (te *testEnv) teardown() error {
	for _, cmd := range te.tearDownCmds {
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

type testLink struct {
	testEnv
	name     string
	linkType string
}

func (tl *testLink) prepTestLink(name, linkType string) error {
	if os.Getuid() != 0 {
		return errors.New("skipping test; must be root")
	}

	tl.name = name
	tl.linkType = linkType

	xpath, err := exec.LookPath("ip")
	if err != nil {
		return err
	}

	tl.createCmds = append(tl.createCmds, &exec.Cmd{
		Path: xpath,
		Args: []string{"ip", "link", "add", tl.name, "type", tl.linkType},
	})

	tl.tearDownCmds = append(tl.tearDownCmds, &exec.Cmd{
		Path: xpath,
		Args: []string{"ip", "link", "del", tl.name},
	})

	return nil
}

func (tl *testLink) prepLinkOpts(opts LinkOptions) error {
	if os.Getuid() != 0 {
		return errors.New("skipping test; must be root")
	}

	macaddr := opts.MacAddr
	mtu := strconv.Itoa(opts.MTU)
	flags := opts.Flags

	xpath, err := exec.LookPath("ip")
	if err != nil {
		return err
	}

	if macaddr != "" {
		tl.setupCmds = append(tl.setupCmds, &exec.Cmd{
			Path: xpath,
			Args: []string{"ip", "link", "set", "dev", tl.name, "address", macaddr},
		})
	}

	if mtu != "" {
		tl.setupCmds = append(tl.setupCmds, &exec.Cmd{
			Path: xpath,
			Args: []string{"ip", "link", "set", "dev", tl.name, "mtu", mtu},
		})
	}

	if (flags & syscall.IFF_UP) == syscall.IFF_UP {
		tl.setupCmds = append(tl.setupCmds, &exec.Cmd{
			Path: xpath,
			Args: []string{"ip", "link", "set", "dev", tl.name, "up"},
		})
	}

	return nil
}

type testDocker struct {
	testEnv
	name    string
	command string
}

func (td *testDocker) prepTestDocker(name, command string) error {
	if os.Getuid() != 0 {
		return errors.New("skipping test; must be root")
	}

	td.name = name
	td.command = command

	xpath, err := exec.LookPath("docker")
	if err != nil {
		return err
	}

	td.createCmds = append(td.createCmds, &exec.Cmd{
		Path: xpath,
		Args: []string{"docker", "run", "-t", "-d", "--name", td.name, "ubuntu", td.command},
	})

	td.tearDownCmds = append(td.tearDownCmds, &exec.Cmd{
		Path: xpath,
		Args: []string{"docker", "stop", td.name},
	})

	td.tearDownCmds = append(td.tearDownCmds, &exec.Cmd{
		Path: xpath,
		Args: []string{"docker", "rm", td.name},
	})

	return nil
}

type testLinkInfo struct {
	linkType string
	linkData string
}

var testLinkInfoData = map[string]func([]string) (string, error){
	"macvlan": macvlanInfo,
	"vlan":    vlanInfo,
	"macvtap": macvtapInfo,
}

func macvlanInfo(data []string) (string, error) {
	if len(data) < 3 {
		return "", fmt.Errorf("Unable to parse macvlan result")
	}

	return data[2], nil
}

func vlanInfo(data []string) (string, error) {
	if len(data) < 5 {
		return "", fmt.Errorf("Unable to parse vlan result")
	}

	return data[4], nil
}

func macvtapInfo(data []string) (string, error) {
	// At the moment this should work fine
	// If we add more tests in the future we will modify this function
	return macvlanInfo(data)
}

func linkInfo(name, linkType string) (*testLinkInfo, error) {
	ipPath, err := exec.LookPath("ip")
	if err != nil {
		return nil, fmt.Errorf("Unable to find ip in PATH: %s", err)
	}

	list := exec.Command(ipPath, "-d", "link", "show", name)
	res := exec.Command("tail", "-1")

	var pipeErr error
	var out bytes.Buffer
	res.Stdin, pipeErr = list.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Unable start UNIX pipe: %s", pipeErr)
	}
	res.Stdout = &out

	if err := res.Start(); err != nil {
		return nil, fmt.Errorf("Unable to tail the result: %s", err)
	}

	if err := list.Run(); err != nil {
		return nil, fmt.Errorf("Unable to retrieve interface information: %s", err)
	}

	if err := res.Wait(); err != nil {
		return nil, fmt.Errorf("Could not read UNIX pipe data: %s", err)
	}

	data := strings.Fields(strings.TrimSpace(out.String()))

	linkInfoFunc, ok := testLinkInfoData[linkType]
	if !ok {
		return &testLinkInfo{
			linkType: data[0],
			linkData: "",
		}, nil
	}

	linkData, err := linkInfoFunc(data)
	if err != nil {
		return nil, fmt.Errorf("Could not read Link Info: %s", err)
	}

	return &testLinkInfo{
		linkType: data[0],
		linkData: linkData,
	}, nil
}
