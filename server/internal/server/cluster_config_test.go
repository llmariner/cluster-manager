package server

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	v1 "github.com/llmariner/cluster-manager/api/v1"
	"github.com/llmariner/cluster-manager/server/internal/config"
	"github.com/llmariner/cluster-manager/server/internal/k8s"
	"github.com/llmariner/cluster-manager/server/internal/store"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestClusterConfig(t *testing.T) {
	st, tearDown := store.NewTest(t)
	defer tearDown()

	fc := fake.NewFakeClient()
	nvidiaConfig := config.NVIDIAConfig{
		DevicePluginConfigMapName:      "llmariner-device-plugin-config",
		DevicePluginConfigMapNamespace: "nvidia",
		DevicePluginConfigName:         "llmariner-default",
	}
	srv := New(st, k8s.NewFakeClientFactory(fc), nvidiaConfig, time.Hour, testr.New(t))
	ctx := fakeAuthInto(context.Background())

	c, err := srv.CreateCluster(ctx, &v1.CreateClusterRequest{
		Name: "cluster",
	})
	assert.NoError(t, err)

	config, err := srv.CreateClusterConfig(ctx, &v1.CreateClusterConfigRequest{
		ClusterId: c.Id,
		DevicePluginConfig: &v1.DevicePluginConfig{
			TimeSlicing: &v1.DevicePluginConfig_TimeSlicing{
				Gpus: 2,
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, int32(2), config.DevicePluginConfig.TimeSlicing.Gpus)

	var cm corev1.ConfigMap
	key := client.ObjectKey{
		Name:      nvidiaConfig.DevicePluginConfigMapName,
		Namespace: nvidiaConfig.DevicePluginConfigMapNamespace,
	}
	err = fc.Get(ctx, key, &cm)
	assert.NoError(t, err)
	_, ok := cm.Data[nvidiaConfig.DevicePluginConfigName]
	assert.True(t, ok)

	config, err = srv.GetClusterConfig(ctx, &v1.GetClusterConfigRequest{
		ClusterId: c.Id,
	})
	assert.NoError(t, err)
	assert.Equal(t, int32(2), config.DevicePluginConfig.TimeSlicing.Gpus)

	_, err = srv.DeleteClusterConfig(ctx, &v1.DeleteClusterConfigRequest{
		ClusterId: c.Id,
	})
	assert.NoError(t, err)

	err = fc.Get(ctx, key, &cm)
	assert.Error(t, err)
	assert.True(t, apierrors.IsNotFound(err))
}

func TestCreateClusterConfig(t *testing.T) {
	dpconfig := &v1.DevicePluginConfig{
		TimeSlicing: &v1.DevicePluginConfig_TimeSlicing{
			Gpus: 2,
		},
	}

	tcs := []struct {
		name    string
		reqF    func(clusterID string) *v1.CreateClusterConfigRequest
		wantErr bool
	}{
		{
			name: "success",
			reqF: func(clusterID string) *v1.CreateClusterConfigRequest {
				return &v1.CreateClusterConfigRequest{
					ClusterId:          clusterID,
					DevicePluginConfig: dpconfig,
				}
			},
		},
		{
			name: "empty cluster ID",
			reqF: func(clusterID string) *v1.CreateClusterConfigRequest {
				return &v1.CreateClusterConfigRequest{
					ClusterId:          "",
					DevicePluginConfig: dpconfig,
				}
			},
			wantErr: true,
		},
		{
			name: "non-existent cluster ID",
			reqF: func(clusterID string) *v1.CreateClusterConfigRequest {
				return &v1.CreateClusterConfigRequest{
					ClusterId:          "non-existent-cluster",
					DevicePluginConfig: dpconfig,
				}
			},
			wantErr: true,
		},
		{
			name: "invalid device plugin config",
			reqF: func(clusterID string) *v1.CreateClusterConfigRequest {
				return &v1.CreateClusterConfigRequest{
					ClusterId: clusterID,
					DevicePluginConfig: &v1.DevicePluginConfig{
						TimeSlicing: &v1.DevicePluginConfig_TimeSlicing{
							Gpus: -1, // Invalid value
						},
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			st, tearDown := store.NewTest(t)
			defer tearDown()

			nvidiaConfig := config.NVIDIAConfig{
				DevicePluginConfigMapName:      "llmariner-device-plugin-config",
				DevicePluginConfigMapNamespace: "nvidia",
				DevicePluginConfigName:         "llmariner-default",
			}
			srv := New(st, k8s.NewFakeClientFactory(fake.NewFakeClient()), nvidiaConfig, time.Hour, testr.New(t))
			ctx := fakeAuthInto(context.Background())

			c, err := srv.CreateCluster(ctx, &v1.CreateClusterRequest{
				Name: "cluster",
			})
			assert.NoError(t, err)

			_, err = srv.CreateClusterConfig(ctx, tc.reqF(c.Id))
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
