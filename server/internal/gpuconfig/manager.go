package gpuconfig

import (
	"context"

	nv1 "github.com/NVIDIA/k8s-device-plugin/api/config/v1"
	"github.com/llmariner/cluster-manager/server/internal/k8s"
	"gopkg.in/yaml.v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// NewManager returns a new Manager.
func NewManager(
	k8sClient k8s.Client,
	configMapName string,
	configMapNamespace string,
	defaultConfigName string,
) *Manager {
	return &Manager{
		k8sClient:          k8sClient,
		configMapName:      configMapName,
		configMapNamespace: configMapNamespace,
		defaultConfigName:  defaultConfigName,
	}
}

// Manager manages the NVIDIA GPU Device Plugin configuration in a Kubernetes cluster.
type Manager struct {
	k8sClient k8s.Client

	configMapName      string
	configMapNamespace string

	// defaultConfigName is the default config name within the ConfigMap for the NVIDIA Device Plugin config.
	defaultConfigName string

	// TODO(kenji): Add clusterPolicyName is the name of the cluster policy
	// (https://github.com/NVIDIA/gpu-operator/blob/main/api/nvidia/v1/clusterpolicy_types.go).
}

// CreateOrUpdateConfigMap creates or updates the ConfigMap for the NVIDIA GPU Device Plugin configuration.
func (m *Manager) CreateOrUpdateConfigMap(ctx context.Context, dpconfig *nv1.Config) error {
	dpConfigBytes, err := yaml.Marshal(dpconfig)
	if err != nil {
		return err
	}

	configData := map[string][]byte{
		m.defaultConfigName: dpConfigBytes,
	}

	if _, err := m.k8sClient.GetConfigMap(ctx, m.configMapName, m.configMapNamespace); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}

		// The ConfigMap does not exist. Create it.
		if err := m.k8sClient.CreateConfigMap(ctx, m.configMapName, m.configMapNamespace, configData); err != nil {
			return err
		}
		return nil
	}

	// The ConfigMap exists. Update it.
	if _, err := m.k8sClient.UpdateConfigMap(ctx, m.configMapName, m.configMapNamespace, configData); err != nil {
		return err
	}

	return nil
}
