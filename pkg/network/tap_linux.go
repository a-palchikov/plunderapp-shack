package network

import (
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

// CreateTap will create a tap device for qemu
func (e *Environment) CreateTap(tapName string) error {
	// Find bridge
	if e.BridgeLink == nil {
		bridge, err := netlink.LinkByName(e.BridgeName)
		if err != nil {
			return err
		}
		e.BridgeLink = bridge
	}

	// Create TAP
	tap := &netlink.Tuntap{LinkAttrs: netlink.LinkAttrs{Name: tapName}, Mode: netlink.TUNTAP_MODE_TAP}
	err := netlink.LinkAdd(tap)
	if err != nil {
		return errors.Wrapf(err, "adding %s", tap.Name)
	}

	// Add Tap to bridge
	err = netlink.LinkSetMaster(tap, e.BridgeLink)
	if err != nil {
		return errors.Wrapf(err, "adding %s to bridge %s", tap.Name, e.BridgeName)
	}
	return nil
}

// DeleteTap will remove a tap device from qemu and the bridge
func (e *Environment) DeleteTap(tapName string) error {
	// Find Tap
	tapLink, err := netlink.LinkByName(tapName)
	if err != nil {
		return errors.Wrapf(err, "deleting %s", tapName)
	}

	// Remove Tap device
	err = netlink.LinkDel(tapLink)
	if err != nil {
		return errors.Wrapf(err, "deleting %s", tapName)
	}
	return nil

}
