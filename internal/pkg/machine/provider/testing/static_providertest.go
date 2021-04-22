package main

import (
	"fmt"

	"github.com/prakashmirji/cluster-api-provider-ecp/internal/pkg/machine"
	"github.com/prakashmirji/cluster-api-provider-ecp/internal/pkg/machine/provider/static"
)

func main() {
	secret := make(map[string]string)
	config := make(map[string]string)
	secret["dbNamespace"] = "static-worker"
	provider := "static"

	macProvider, err := newMachineProvider(provider, secret, config)
	if err != nil {
		fmt.Errorf("unable to create machine provider client, %v\n", err)
		//r.setMachineError(kubeMachine, infrav1alpha2.MachineProviderError, msg)
		return
	}

	machineName := "worker1"
	tags := make(map[string]string)

	mac, err := macProvider.GetMachine(machineName)
	if err != nil {
		fmt.Errorf("error in getting machine, error: %v\n", err)
	}
	if mac == nil {
		var networks []string
		createInfo := machine.MachineCreateData{
			OsImage:   "",
			OsVersion: "",
			Size:      "",
			SSHUser:   "test",
			SSHKey:    "test",
			Networks:  networks,
			Proxy:     "",
			Tags:      tags,
		}
		fmt.Printf("machine create started\n")
		mac, err := macProvider.CreateMachine(machineName, createInfo)
		if err != nil {
			fmt.Errorf("unable to create machine, error: %v\n", err)
			return
		}
		fmt.Printf("successfully got machine: %s, ID: %v, hostname: %s \n", mac.Name(), mac.ID(), mac.Hostname())
	} else {
		fmt.Printf("successfully got machine: %s, ID: %v, hostname: %s \n", mac.Name(), mac.ID(), mac.Hostname())
	}
}

func newMachineProvider(provider string, providerSecret, providerConfig map[string]string) (machine.Provider, error) {
	switch provider {
	case "static":
		return static.NewStaticProvider(providerSecret, providerConfig)
	default:
		return nil, fmt.Errorf("machine provider '%v' not implemented", provider)
	}
}
