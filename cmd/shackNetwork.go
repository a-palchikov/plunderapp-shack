package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/plunder-app/shack/pkg/network"
	"github.com/spf13/cobra"
)

func init() {
	shackNetwork.AddCommand(shackNetworkCreate)
	shackNetwork.AddCommand(shackNetworkCheck)
	shackNetwork.AddCommand(shackNetworkDelete)
	shackNetwork.AddCommand(shackNetworkNat)
}

var shackNetworkCheck = &cobra.Command{
	Use:   "check",
	Short: "Check the bridge",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("shack Networking configuration")
		cfg, err := network.OpenFile(configPath)
		if err != nil {
			return errors.Wrap(err, "reading configuration")
		}

		err = cfg.CheckBridge()
		if err != nil {
			return errors.Wrap(err, "validating the bridge")
		}
		return nil
	},
}

var shackNetworkCreate = &cobra.Command{
	Use:   "create",
	Short: "Create the bridge",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("shack Networking configuration")
		cfg, err := network.OpenFile(configPath)
		if err != nil {
			return errors.Wrap(err, "reading configuration")
		}

		err = cfg.CreateBridge()
		if err != nil {
			return errors.Wrap(err, "creating the bridge")
		}

		err = cfg.AddBridgeAddress()
		if err != nil {
			return errors.Wrap(err, "creating the bridge")
		}

		err = cfg.BridgeUp()
		if err != nil {
			return errors.Wrap(err, "creating the bridge")
		}
		return nil
	},
}

var shackNetworkDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete the bridge",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("shack Networking configuration\n")
		cfg, err := network.OpenFile(configPath)
		if err != nil {
			return errors.Wrap(err, "reading configuration")
		}

		err = cfg.DeleteBridge()
		if err != nil {
			return errors.Wrap(err, "deleting the bridge")
		}
		return nil
	},
}

var shackNetworkNat = &cobra.Command{
	Use:   "nat",
	Short: "Enable NAT",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("shack Networking configuration")
		cfg, err := network.OpenFile(configPath)
		if err != nil {
			return errors.Wrap(err, "reading configuration")
		}

		err = cfg.EnableNat()
		if err != nil {
			return errors.Wrap(err, "enabling NAT")
		}
		return nil
	},
}
