// (c) Copyright 2019-2020 Hewlett Packard Enterprise Development LP

// Package static provides an implementation of a static machine provider based on config maps
package static

import (
	"fmt"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/prakashmirji/cluster-api-provider-ecp/internal/pkg/machine"
)

func getLogger() logr.Logger {
	return ctrl.Log.WithName("machineProvider").WithName("static")
}

// Ensure that StaticProvider is a Provider.
var _ machine.Provider = (*StaticProvider)(nil)

// Ensure that StaticProviderMachine is a Machine.
var _ machine.Machine = (*StaticProviderMachine)(nil)

// StaticProviderMachine wraps a staticProvider machine in a generic Machine implementation.
type StaticProviderMachine struct {
	id       string
	name     string
	hostname string
	ipaddr   string
}

// ID returns Machine ID
func (m *StaticProviderMachine) ID() string {
	return m.id
}

// Name returns Machine Name
func (m *StaticProviderMachine) Name() string {
	return m.name
}

// Hostname returns Machine Hostname
func (m *StaticProviderMachine) Hostname() string {
	return m.hostname
}

// IPAddr returns Machine ip address
func (m *StaticProviderMachine) IPAddr() string {
	return m.ipaddr
}

// StaticProvider represents a StaticProvider client struct
type StaticProvider struct {
	clientset   *kubernetes.Clientset
	dbNamespace string
	logger      logr.Logger
}

// GetInfraProviderStatus ...
// func (mp *StaticProvider) GetInfraProviderStatus() (*v1alpha2.KubeInfraProviderStatus, error) {
// 	return &v1alpha2.KubeInfraProviderStatus{
// 		Ready: true,
// 	}, nil
// }

// CreateMachine creates a machine
func (mp *StaticProvider) CreateMachine(name string, data machine.MachineCreateData) (machine.Machine, error) {
	dbClient := mp.clientset.CoreV1().ConfigMaps(mp.dbNamespace)

	// Fetch an available machine to be allocated
	availableList, err := dbClient.List(metav1.ListOptions{LabelSelector: "allocated=false"})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch available machine config maps: %v", err)
	}
	if len(availableList.Items) == 0 {
		return nil, fmt.Errorf("no more machines available to allocate")
	}
	allocMac := availableList.Items[0]

	// update the allocated machine details
	allocMac.Labels["allocated"] = "true"
	allocMac.Labels["name"] = name
	_, err = dbClient.Update(&allocMac)
	if err != nil {
		return nil, fmt.Errorf("failed to update machine config map: %v", err)
	}

	fmt.Printf("updated machine config map for name: %s\n", name)

	return &StaticProviderMachine{
		id:       allocMac.Name,
		name:     name,
		hostname: allocMac.Data["hostname"],
		ipaddr:   allocMac.Data["ipAddr"],
	}, nil
}

// GetMachine fetches a machine
func (mp *StaticProvider) GetMachine(name string) (machine.Machine, error) {
	dbClient := mp.clientset.CoreV1().ConfigMaps(mp.dbNamespace)

	allocatedMachine, err := dbClient.List(metav1.ListOptions{LabelSelector: fmt.Sprintf("name=%s", name)})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch allocated machine config maps: %v", err)
	}
	if len(allocatedMachine.Items) == 0 {
		return nil, nil
	}
	// if len(allocatedMachine.Items) > 1 {
	// 	return nil, fmt.Errorf("invalid number of machines returned for name %q: %d", name, len(allocatedMachine.Items))
	// }
	return &StaticProviderMachine{
		id:       allocatedMachine.Items[0].Name,
		name:     name,
		hostname: allocatedMachine.Items[0].Data["hostname"],
		ipaddr:   allocatedMachine.Items[0].Data["ipAddr"],
	}, nil
}

// DeleteMachine deletes a machine
func (mp *StaticProvider) DeleteMachine(name string) error {
	dbClient := mp.clientset.CoreV1().ConfigMaps(mp.dbNamespace)

	allocatedMachine, err := dbClient.List(metav1.ListOptions{LabelSelector: fmt.Sprintf("name=%s", name)})
	if err != nil {
		return fmt.Errorf("failed to fetch allocated machine config maps: %v", err)
	}
	if len(allocatedMachine.Items) != 1 {
		return fmt.Errorf("invalid number of machines returned for name %q: %v", name, err)
	}
	allocMac := allocatedMachine.Items[0]

	// update the details of the machine
	allocMac.Labels["allocated"] = "false"
	allocMac.Labels["name"] = ""

	_, err = dbClient.Update(&allocMac)
	if err != nil {
		return fmt.Errorf("failed to update machine config map: %v", err)
	}
	return nil
}

// NewStaticProvider returns a Static implementation of a machine Provider
func NewStaticProvider(providerSecret, providerConfig map[string]string) (machine.Provider, error) {
	var dbNamespace string
	if val, ok := providerSecret["dbNamespace"]; ok {
		dbNamespace = val
	}
	// creates the in-cluster config
	//config, err := rest.InClusterConfig()

	var kubeconfig *string
	filename := "/home/pmirji/.kube/config"
	kubeconfig = &filename
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize static provider: %v", err)
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize static provider clientset: %v", err)
	}
	return &StaticProvider{clientset: clientset, dbNamespace: dbNamespace, logger: getLogger()}, nil
}
