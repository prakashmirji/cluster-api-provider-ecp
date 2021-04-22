// (c) Copyright 2019-2021 Hewlett Packard Enterprise Development LP

package machine

// Provider defines the interface that Cassini uses to control k8s
// machine providers.
type Provider interface {
	//GetInfraProviderStatus() (*v1alpha3.KubeInfraProviderStatus, error)
	CreateMachine(string, MachineCreateData) (Machine, error)
	GetMachine(string) (Machine, error)
	DeleteMachine(string) error
}

// Machine defines the interface that Cassini uses to represent k8s
// machine.
type Machine interface {
	ID() string
	Name() string
	Hostname() string
	IPAddr() string
}

// MachineCreateData holds the generic parameters to create a machine
type MachineCreateData struct {
	OsImage   string            `json:"osImage,omitempty"`   // OS flavor for the machine
	OsVersion string            `json:"osVersion,omitempty"` // OS version for the machine
	Size      string            `json:"size,omitempty"`      // ProvisioningPlan of the machine
	SSHKey    string            `json:"sshKey,omitempty"`    // name of the SSH-Key to be used
	SSHUser   string            `json:"sshUser,omitempty"`   // name of the SSH-User to be used
	Networks  []string          `json:"networks,omitempty"`  // List of names of networks to be attached to the machine
	Proxy     string            `json:"proxy,omitempty"`     // Network proxy address for external communication
	Tags      map[string]string `json:"tags,omitempty"`      //Tags constains a list of tags to be passed while machine creation`
}

// State represents the state of the entity
type State int

const (
	// StateUnknown represents an unknown state
	StateUnknown State = iota
	// StateInfraProvisioning represents infra being provisioned
	StateInfraProvisioning
	// StateCreating represents entity being created
	StateCreating
	// StateReady represents ready entity
	StateReady
	// StateUpdating represents entity being updated
	StateUpdating
	// StateUpgrading represents entity being upgraded
	StateUpgrading
	// StateDeleting represents entity being deleted
	StateDeleting
	// StateDeleted represent deleted entity
	StateDeleted
	// StateError represent errored entity
	StateError
)

// String gives the string representation of State
func (s State) String() string {
	switch s {
	case StateInfraProvisioning:
		return "infra-provisioning"
	case StateCreating:
		return "creating"
	case StateReady:
		return "ready"
	case StateUpdating:
		return "updating"
	case StateUpgrading:
		return "upgrading"
	case StateDeleting:
		return "deleting"
	case StateDeleted:
		return "deleted"
	case StateError:
		return "error"
	}
	return "unknown"
}
