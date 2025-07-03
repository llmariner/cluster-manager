package gpuconfig

// The following definitions are twekaed verisons of github.com/Nvidia/k8s-device-plugin@v0.17.2/api/config/v1.
// 'omitempty' is added to some of the fields.

// Config is a versioned struct used to hold configuration information.
type Config struct {
	Version string  `yaml:"version"`
	Flags   Flags   `yaml:"flags,omitempty"`
	Sharing Sharing `yaml:"sharing,omitempty"`
}

// Flags contains the command line flags used.
type Flags struct {
	CommandLineFlags
}

// CommandLineFlags holds the list of command line flags used to configure the device plugin and GFD.
type CommandLineFlags struct {
	MigStrategy             *string `yaml:"migStrategy"`
	FailOnInitError         *bool   `yaml:"failOnInitError,omitempty"`
	MpsRoot                 *string `yaml:"mpsRoot,omitempty"`
	NvidiaDriverRoot        *string `yaml:"nvidiaDriverRoot,omitempty"`
	NvidiaDevRoot           *string `yaml:"nvidiaDevRoot,omitempty"`
	GDSEnabled              *bool   `yaml:"gdsEnabled,omitempty"`
	MOFEDEnabled            *bool   `yaml:"mofedEnabled,omitempty"`
	UseNodeFeatureAPI       *bool   `yaml:"useNodeFeatureAPI,omitempty"`
	DeviceDiscoveryStrategy *string `yaml:"deviceDiscoveryStrategy,omitempty"`
}

// Sharing encapsulates the set of sharing strategies that are supported.
type Sharing struct {
	// TimeSlicing defines the set of replicas to be made for timeSlicing available resources.
	TimeSlicing ReplicatedResources `yaml:"timeSlicing,omitempty"`
	// MPS defines the set of replicas to be shared using MPS
	MPS *ReplicatedResources `yaml:"mps,omitempty"`
}

// ReplicatedResources defines generic options for replicating devices.
type ReplicatedResources struct {
	RenameByDefault            bool                 `yaml:"renameByDefault,omitempty"`
	FailRequestsGreaterThanOne bool                 `yaml:"failRequestsGreaterThanOne,omitempty"`
	Resources                  []ReplicatedResource `yaml:"resources,omitempty"`
}

// ReplicatedResource represents a resource to be replicated.
type ReplicatedResource struct {
	Name     ResourceName      `yaml:"name"`
	Rename   ResourceName      `yaml:"rename,omitempty"`
	Devices  ReplicatedDevices `yaml:"devices,flow,omitempty"`
	Replicas int               `yaml:"replicas"`
}

// ResourceName represents a valid resource name in Kubernetes
type ResourceName string

// ReplicatedDevices encapsulates the set of devices that should be replicated for a given resource.
// This struct should be treated as a 'union' and only one of the fields in this struct should be set at any given time.
type ReplicatedDevices struct {
	All   bool
	Count int
	List  []ReplicatedDeviceRef
}

// ReplicatedDeviceRef can either be a full GPU index, a MIG index, or a UUID (full GPU or MIG)
type ReplicatedDeviceRef string

// CreateTimeSlicingDevicePluginConfig returns a v1.Config that configures the NVIDIA GPU device plugin for time slicing.
//
// Link: https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/latest/gpu-sharing.html
func CreateTimeSlicingDevicePluginConfig(gpus int) *Config {
	return &Config{
		Version: "v1",
		Flags: Flags{
			CommandLineFlags: CommandLineFlags{
				MigStrategy: strPtr("none"),
			},
		},
		Sharing: Sharing{
			TimeSlicing: ReplicatedResources{
				// Keep the original resource name. When set to true, each resource is advertised under
				// the name <resource-name>.shared instead of <resource-name>.
				RenameByDefault: false,
				// This is to enforce awareness that requesting more than one GPU replica does not result
				// in receiving more proportional access to the GPU.
				//
				// For example, if 4 GPU replicas are available and two pods request 1 GPU each and
				// a third pod requests 2 GPUs, the applications in the three pods have an equal share of GPU
				// compute time. Specifically, the pod that requests 2 GPUs does not receive twice as much compute time
				// as the pods that request 1 GPU.
				//
				// When set to true, a resource request for more than one GPU fails with an UnexpectedAdmissionError.
				// In this case, you must manually delete the pod, update the resource request, and redeploy.
				FailRequestsGreaterThanOne: true,
				Resources: []ReplicatedResource{
					{
						Name:     "nvidia.com/gpu",
						Replicas: gpus,
					},
				},
			},
		},
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
