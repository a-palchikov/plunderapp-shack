package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/plunder-app/shack/pkg/network"
	"github.com/plunder-app/shack/pkg/vmm"
)

var vmUUID string
var foreground, vnc, disk bool

func init() {
	shackVM.PersistentFlags().StringVar(&vmUUID, "id", "000000", "The UUID for a virtual machine")
	shackVMStart.Flags().BoolVarP(&foreground, "foreground", "f", false, "Whether to start VM in foreground")
	shackVMStart.Flags().BoolVarP(&vnc, "vnc", "v", false, "Enable VNC")
	shackVMStop.Flags().BoolVarP(&disk, "disk", "d", false, "Delete Disk")

	// Add subcommands
	shackVM.AddCommand(shackVMStart)
	shackVM.AddCommand(shackVMStop)
}

var shackVM = &cobra.Command{
	Use:   "vm",
	Short: "Manage VMs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("shack VM configuration")
		_ = cmd.Help()
	},
}

var shackVMStart = &cobra.Command{
	Use:   "start",
	Short: "Start a Virtual Machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("shack VM configuration")
		cfg, err := network.OpenFile(configPath)
		if err != nil {
			return errors.Wrap(err, "reading configuration")
		}
		if vmUUID == "000000" {
			// Generate VM UUID
			b, err := vmm.GenVMUUID()
			if err != nil {
				return errors.Wrap(err, "generating a VM UID")
			}
			vmUUID = fmt.Sprintf("%02x%02x%02x", b[0], b[1], b[2])
		}

		vmInterface := fmt.Sprintf("%s-%s", cfg.NicPrefix, vmUUID)
		if len(vmInterface) > 15 {
			return errors.Errorf("The interface name [%s] is too long for the interface standard, shorten the nicPrefix", vmInterface)
		}
		// Create Tap Device (and add to bridge)
		if err := cfg.CreateTap(vmInterface); err != nil {
			return errors.Wrapf(err, "creating TAP for %s", vmInterface)
		}

		// Generate MAC address using UUID and Mac prefix
		mac := vmm.GenVMMac(cfg.NicMacPrefix, vmUUID)

		// Create Disk
		err = vmm.CreateDisk(vmUUID, "4G")
		if err != nil {
			return errors.Wrap(err, "creating a disk")
		}

		// Start Virtual Machine
		err = vmm.Start(mac, vmUUID, cfg.NicPrefix, foreground, vnc)
		if err != nil {
			return errors.Wrapf(err, "starting VM %s", vmUUID)
		}

		// If this is ran in the foreground then we will want to tidy up the created interface
		if foreground {
			log.Infof("Deleting interface [%s]", vmInterface)
			err = cfg.DeleteTap(vmInterface)
			if err != nil {
				return errors.Wrap(err, "removing the TAP device")
			}
		}
		return nil
	},
}

var shackVMStop = &cobra.Command{
	Use:   "stop",
	Short: "Stop a Virtual Machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := network.OpenFile(configPath)
		if err != nil {
			return errors.Wrap(err, "reading configuration")
		}

		// Stop Virtual Machine
		err = vmm.Stop(vmUUID)
		if err != nil {
			return errors.Wrapf(err, "stopping VM %s", vmUUID)
		}

		// Remove Networking configuration
		err = cfg.DeleteTap(fmt.Sprintf("%s-%s", cfg.NicPrefix, vmUUID))
		if err != nil {
			return errors.Wrap(err, "deleting the TAP device")
		}

		if !disk {
			return nil
		}

		// Delete the disk
		err = vmm.DeleteDisk(vmUUID)
		if err != nil {
			return errors.Wrap(err, "deleting the disk")
		}
		return nil
	},
}
