package gpuconfig

import (
	"context"
	"testing"

	nv1 "github.com/NVIDIA/k8s-device-plugin/api/config/v1"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestClient(t *testing.T) {
	fc := fake.NewFakeClient()

	configClient := NewClient(
		newFakeK8sClient(fc),
		"cmName",
		"cmNamespace",
		"gpuConfigName",
	)

	ctx := context.Background()
	dpconfig := CreateTimeSlicingDevicePluginConfig(2)
	err := configClient.CreateOrUpdateConfigMap(ctx, dpconfig)
	assert.NoError(t, err)

	var cm corev1.ConfigMap
	err = fc.Get(ctx, client.ObjectKey{Name: "cmName", Namespace: "cmNamespace"}, &cm)
	assert.NoError(t, err)

	configData := cm.Data["gpuConfigName"]
	var got nv1.Config
	err = yaml.Unmarshal([]byte(configData), &got)
	assert.NoError(t, err)

	res := got.Sharing.TimeSlicing.Resources[0]
	assert.Equal(t, "nvidia.com/gpu", string(res.Name))
	assert.Equal(t, 2, res.Replicas)
}

func newFakeK8sClient(c client.WithWatch) *fakeK8sClient {
	return &fakeK8sClient{c: c}
}

type fakeK8sClient struct {
	c client.WithWatch
}

func (c *fakeK8sClient) GetConfigMap(ctx context.Context, name, namespace string) (*corev1.ConfigMap, error) {
	var cm corev1.ConfigMap
	if err := c.c.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &cm); err != nil {
		return nil, err
	}
	return &cm, nil
}

func (c *fakeK8sClient) CreateConfigMap(ctx context.Context, name, namespace string, data map[string][]byte) error {
	cm := configMap(name, namespace, data)
	if err := c.c.Create(ctx, cm); err != nil {
		return err
	}

	return nil
}

func (c *fakeK8sClient) UpdateConfigMap(ctx context.Context, name, namespace string, data map[string][]byte) (*corev1.ConfigMap, error) {
	cm := configMap(name, namespace, data)
	if err := c.c.Update(ctx, cm); err != nil {
		return nil, err
	}
	return cm, nil
}

func configMap(name, namespace string, data map[string][]byte) *corev1.ConfigMap {
	d := map[string]string{}
	for k, v := range data {
		d[k] = string(v)
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: d,
	}
}
