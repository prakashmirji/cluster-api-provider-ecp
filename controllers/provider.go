package controllers

import (
	"fmt"

	"github.com/prakashmirji/cluster-api-provider-ecp/internal/pkg/machine"
	"github.com/prakashmirji/cluster-api-provider-ecp/internal/pkg/machine/provider/static"
)

func newMachineProvider(provider string, providerSecret, providerConfig map[string]string) (machine.Provider, error) {
	switch provider {
	case "static":
		return static.NewStaticProvider(providerSecret, providerConfig)
	default:
		return nil, fmt.Errorf("machine provider '%v' not implemented", provider)
	}
}
