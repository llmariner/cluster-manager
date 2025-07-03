package gpuconfig

import (
	"context"

	"github.com/llmariner/cluster-manager/server/internal/k8s"
	"gopkg.in/yaml.v2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// NewClient returns a new Client.
func NewClient(
	k8sClient k8s.Client,
	configMapName string,
	configMapNamespace string,
	defaultConfigName string,
) *Client {
	return &Client{
		k8sClient:          k8sClient,
		configMapName:      configMapName,
		configMapNamespace: configMapNamespace,
		defaultConfigName:  defaultConfigName,
	}
}

// Client manages the NVIDIA GPU Device Plugin configuration in a Kubernetes cluster.
type Client struct {
	k8sClient k8s.Client

	configMapName      string
	configMapNamespace string

	// defaultConfigName is the default config name within the ConfigMap for the NVIDIA Device Plugin config.
	defaultConfigName string

	// TODO(kenji): Add clusterPolicyName is the name of the cluster policy
	// (https://github.com/NVIDIA/gpu-operator/blob/main/api/nvidia/v1/clusterpolicy_types.go).
}

// CreateOrUpdateConfigMap creates or updates the ConfigMap for the NVIDIA GPU Device Plugin configuration.
func (c *Client) CreateOrUpdateConfigMap(ctx context.Context, dpconfig *Config) error {
	dpConfigBytes, err := yaml.Marshal(dpconfig)
	if err != nil {
		return err
	}

	configData := map[string]string{
		c.defaultConfigName: string(dpConfigBytes),
	}

	if _, err := c.k8sClient.GetConfigMap(ctx, c.configMapName, c.configMapNamespace); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}

		// The ConfigMap does not exist. Create it.
		if err := c.k8sClient.CreateConfigMap(ctx, c.configMapName, c.configMapNamespace, configData); err != nil {
			return err
		}
		return nil
	}

	// The ConfigMap exists. Update it.
	if _, err := c.k8sClient.UpdateConfigMap(ctx, c.configMapName, c.configMapNamespace, configData); err != nil {
		return err
	}

	return nil
}

// DeleteConfigMapIfExists deletes the ConfigMap for the NVIDIA GPU Device Plugin configuration if it exists.
func (c *Client) DeleteConfigMapIfExists(ctx context.Context) error {
	if _, err := c.k8sClient.GetConfigMap(ctx, c.configMapName, c.configMapNamespace); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if err := c.k8sClient.DeleteConfigMap(ctx, c.configMapName, c.configMapNamespace); err != nil {
		return err
	}
	return nil
}
