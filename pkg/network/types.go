package network

import (
	"os"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

// Environment defines the configuration of the shack environment
type Environment struct {
	// Host configuration
	Interface string `json:"interface"`

	// Bridge configuration
	BridgeName    string `json:"bridgeName"`
	BridgeAddress string `json:"bridgeAddress"`

	// Used during runtime
	BridgeLink netlink.Link `json:"-"`

	// VM Nic configuration
	NicPrefix    string `json:"nicPrefix"`
	NicMacPrefix string `json:"nicMacPrefix"`
}

// OpenFile will open an file and parse the contents
func OpenFile(path string) (*Environment, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.Errorf("finding file at %q", path)
	}

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "reading file")
	}

	var environ Environment
	err = yaml.Unmarshal(yamlFile, &environ)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling YAML file")
	}

	return &environ, nil
}

// ExampleConfig will return a config output
func ExampleConfig() string {
	cfg := Environment{
		Interface:     "eth0",
		BridgeAddress: "192.168.1.1/24",
		BridgeName:    "plunder",
		NicPrefix:     "plndrVM",
		NicMacPrefix:  "c0:ff:ee:",
	}

	b, _ := yaml.Marshal(cfg)
	return string(b)
}
