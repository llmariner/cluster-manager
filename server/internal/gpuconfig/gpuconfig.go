package gpuconfig

import (
	v1 "github.com/NVIDIA/k8s-device-plugin/api/config/v1"
)

// Doc:

// timeSlicing returns a v1.Config that configures the NVIDIA GPU device plugin for time slicing.
//
// Link: https://docs.nvidia.com/datacenter/cloud-native/gpu-operator/latest/gpu-sharing.html
func timeSlicingDevicePluginConfig(gpus int) *v1.Config {
	return &v1.Config{
		Flags: v1.Flags{
			CommandLineFlags: v1.CommandLineFlags{
				MigStrategy: strPtr(v1.MigStrategyNone),
			},
		},
		Sharing: v1.Sharing{
			TimeSlicing: v1.ReplicatedResources{
				// Keep the original resoure name. When set to true, each resource is advertised under
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
				Resources: []v1.ReplicatedResource{
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
