package gpuconfig

import (
	"context"
	"testing"

	nv1 "github.com/NVIDIA/k8s-device-plugin/api/config/v1"
	"github.com/llmariner/cluster-manager/server/internal/k8s"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestClient(t *testing.T) {
	fc := fake.NewFakeClient()

	configClient := NewClient(
		k8s.NewFakeClient(fc),
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

	err = configClient.DeleteConfigMapIfExists(ctx, "cmName", "cmNamespace")
	assert.NoError(t, err)

	err = fc.Get(ctx, client.ObjectKey{Name: "cmName", Namespace: "cmNamespace"}, &cm)
	assert.Error(t, err)
	assert.True(t, apierrors.IsNotFound(err))

	err = configClient.DeleteConfigMapIfExists(ctx, "cmName", "cmNamespace")
	assert.NoError(t, err)

}
