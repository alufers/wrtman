package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type OpenWrtConnection struct {
	HasDHCP             bool
	SSHAddress          string
	SSHClientConfig     *ssh.ClientConfig
	Hostname            string
	sshConn             *ssh.Client
	disconnectorRunning bool
	disconnectTicks     int
	connMutex           sync.Mutex
	WirelessDataService *WirelessDataService
}

func NewOpenWrtConnection(SSHAddress string, clientConfig *ssh.ClientConfig) (conn *OpenWrtConnection) {
	conn = &OpenWrtConnection{
		SSHAddress:      SSHAddress,
		SSHClientConfig: clientConfig,
		connMutex:       sync.Mutex{},
	}
	conn.WirelessDataService = NewWirelessDataService(conn)
	return
}

func (o *OpenWrtConnection) ConnectToSSH() error {
	o.connMutex.Lock()
	defer o.connMutex.Unlock()

	conn, err := ssh.Dial("tcp", o.SSHAddress, o.SSHClientConfig)
	if err != nil {
		return fmt.Errorf("failed to dial %v@%v: %w", o.SSHClientConfig.User, o.SSHAddress, err)
	}

	o.sshConn = conn
	o.disconnectTicks = 0
	if !o.disconnectorRunning {
		go o.disconnectorRoutine()
	}
	// if o.Hostname == "" {
	// 	err := o.DiscoverHostname()
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (o *OpenWrtConnection) ConnectOrGetConn() (*ssh.Client, error) {
	if o.sshConn == nil {
		if err := o.ConnectToSSH(); err != nil {
			return nil, err
		}

	} else {
		o.connMutex.Lock()
		defer o.connMutex.Unlock()
		o.disconnectTicks = 0
	}
	return o.sshConn, nil
}

func (o *OpenWrtConnection) DiscoverHostname() error {

	hn, err := o.RunCommandAndGetString("cat /proc/sys/kernel/hostname")
	if err != nil {
		return fmt.Errorf("failed to discover hostname of %v: %w", o.SSHAddress, err)
	}

	o.Hostname = strings.TrimSpace(hn)
	return nil
}

func (o *OpenWrtConnection) GetUptime() (time.Duration, error) {

	uptimeStr, err := o.RunCommandAndGetString("cat /proc/uptime")
	if err != nil {
		return 0, fmt.Errorf("failed to get uptime of %v: %w", o.SSHAddress, err)
	}
	var uptimeF float64
	fmt.Sscanf(uptimeStr, "%v", &uptimeF)

	return time.Second * time.Duration(uptimeF), nil
}

func (o *OpenWrtConnection) GetMacAddrs() ([]string, error) {
	macsStr, err := o.RunCommandAndGetString("cat /sys/class/net/*/address")
	if err != nil {
		return nil, fmt.Errorf("failed to get mac addrs of %v: %w", o.SSHAddress, err)
	}
	segments := strings.Split(macsStr, "\n")
	out := []string{}
	for _, seg := range segments {
		tr := strings.TrimSpace(seg)
		if tr != "" && tr != "00:00:00:00:00:00" {
			out = append(out, tr)
		}
	}

	return out, nil
}

func (o *OpenWrtConnection) UbusCall(service, method string, out interface{}) error {
	ubusOutput, err := o.RunCommandAndGetString(fmt.Sprintf("ubus call %v %v", service, method))
	if err != nil {
		return fmt.Errorf("failed make an ubus call to %v %v: %w", service, method, err)
	}

	return json.Unmarshal([]byte(ubusOutput), out)
}

func (o *OpenWrtConnection) RunCommandAndGetString(cmd string) (string, error) {
	// log.Printf("run command start at %v: %v", o.SSHAddress, cmd)
	conn, err := o.ConnectOrGetConn()
	if err != nil {
		return "", err
	}

	// log.Printf("run command connection oks: %v", cmd)

	sess, err := conn.NewSession()

	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			o.sshConn = nil
			return o.RunCommandAndGetString(cmd)
		}
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer sess.Close()
	buf := &bytes.Buffer{}
	sess.Stdout = buf
	// sess.Stderr = buf
	err = sess.Run(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to execute '%v': %w", cmd, err)
	}
	// log.Printf("run command done: %v, %v", cmd, buf.String())
	return buf.String(), err
}

func (o *OpenWrtConnection) disconnectorRoutine() {
	o.disconnectorRunning = true
	o.disconnectTicks = 0

	defer func() {
		o.sshConn.Close()
		o.sshConn = nil
		o.disconnectTicks = 0
		o.disconnectorRunning = false
	}()
	for {
		shouldLeave := false

		func() {

			time.Sleep(time.Second * 30)
			o.connMutex.Lock()
			defer o.connMutex.Unlock()
			o.disconnectTicks++
			// log.Printf("disconnectTicks = %v", o.disconnectTicks)
			if o.disconnectTicks > 20 { // 10 minutes
				shouldLeave = true
			}
		}()

		if shouldLeave {
			log.Printf("Disconnection from %v due to inactivity", o.Hostname)
			break
		}

	}
}
