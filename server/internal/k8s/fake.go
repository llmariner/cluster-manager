package k8s

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewFakeClientFactory creates a new fake Kubernetes client factory.
func NewFakeClientFactory(c client.WithWatch) *FakeClientFactory {
	return &FakeClientFactory{c: c}
}

// FakeClientFactory is a fake Kubernetes client factory.
type FakeClientFactory struct {
	c client.WithWatch
}

// NewClient creates a new fake Kubernetes client that implements the Client interface.
func (f *FakeClientFactory) NewClient(clusterID, token string) (Client, error) {
	return NewFakeClient(f.c), nil
}

// NewDynamicClient creates a new fake Kubernetes dynamic client that implements the DynamicClient interface.
func (f *FakeClientFactory) NewDynamicClient(clusterID, token string) (DynamicClient, error) {
	return nil, fmt.Errorf("not implemented")
}

// NewFakeClient creates a new fake Kubernetes client that implements the Client interface.
func NewFakeClient(c client.WithWatch) *FakeK8sClient {
	return &FakeK8sClient{c: c}
}

// FakeK8sClient is a fake Kubernetes client that implements the Client interface.
type FakeK8sClient struct {
	c client.WithWatch
}

// CreateConfigMap creates a new ConfigMap.
func (c *FakeK8sClient) CreateConfigMap(ctx context.Context, name, namespace string, data map[string]string) error {
	cm := configMap(name, namespace, data)
	if err := c.c.Create(ctx, cm); err != nil {
		return err
	}

	return nil
}

// GetConfigMap retrieves a ConfigMap.
func (c *FakeK8sClient) GetConfigMap(ctx context.Context, name, namespace string) (*corev1.ConfigMap, error) {
	var cm corev1.ConfigMap
	if err := c.c.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &cm); err != nil {
		return nil, err
	}
	return &cm, nil
}

// UpdateConfigMap updates an existing ConfigMap.
func (c *FakeK8sClient) UpdateConfigMap(ctx context.Context, name, namespace string, data map[string]string) (*corev1.ConfigMap, error) {
	cm := configMap(name, namespace, data)
	if err := c.c.Update(ctx, cm); err != nil {
		return nil, err
	}
	return cm, nil
}

// DeleteConfigMap deletes a ConfigMap.
func (c *FakeK8sClient) DeleteConfigMap(ctx context.Context, name, namespace string) error {
	var cm corev1.ConfigMap
	if err := c.c.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, &cm); err != nil {
		return err
	}
	if err := c.c.Delete(ctx, &cm); err != nil {
		return err
	}
	return nil
}

func configMap(name, namespace string, data map[string]string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
}
