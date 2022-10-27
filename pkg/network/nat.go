package network

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// EnableNat will configure the kernel to enable natting
func (e *Environment) EnableNat() error {
	file, err := os.OpenFile("/proc/sys/net/ipv4/ip_forward", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return errors.Wrap(err, "opening ip_forward file")
	}
	defer file.Close()
	_, err = file.WriteString("1")
	if err != nil {
		return errors.Wrap(err, "writing ip_forward file")
	}

	if _, err := exec.Command("iptables", "-A", "FORWARD", "-i", e.BridgeName, "-o", e.Interface, "-j", "ACCEPT").CombinedOutput(); err != nil {
		return errors.WithStack(err)
	}

	if _, err := exec.Command("iptables", "-A", "FORWARD", "-i", e.Interface, "-o", e.BridgeName, "-m", "state", "--state", "ESTABLISHED,RELATED", "-j", "ACCEPT").CombinedOutput(); err != nil {
		return errors.WithStack(err)
	}

	if _, err := exec.Command("iptables", "-t", "nat", "-A", "POSTROUTING", "-o", e.Interface, "-j", "MASQUERADE").CombinedOutput(); err != nil {
		return errors.WithStack(err)
	}

	//"iptables", "-A", "FORWARD", "-i", "ens192", "-o", "ens160", "-j", "ACCEPT"
	//"iptables", "-A", "FORWARD", "-i", "ens160", "-o", "ens192", "-m", "state", "--state", "ESTABLISHED,RELATED", "-j", "ACCEPT"
	//"iptables", "-t", "nat", "-A", "POSTROUTING", "-o", "ens160", "-j", "MASQUERADE"

	return nil
}
