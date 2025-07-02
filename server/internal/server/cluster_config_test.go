package server

import (
	"context"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	v1 "github.com/llmariner/cluster-manager/api/v1"
	"github.com/llmariner/cluster-manager/server/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestClusterConfig(t *testing.T) {
	st, tearDown := store.NewTest(t)
	defer tearDown()

	srv := New(st, testr.New(t), time.Hour)
	ctx := fakeAuthInto(context.Background())

	c, err := srv.CreateCluster(ctx, &v1.CreateClusterRequest{
		Name: "cluster",
	})
	assert.NoError(t, err)

	_, err = srv.CreateClusterConfig(ctx, &v1.CreateClusterConfigRequest{
		ClusterId: c.Id,
		DevicePluginConfig: &v1.DevicePluginConfig{
			TimeSlicing: &v1.DevicePluginConfig_TimeSlicing{
				Gpus: 2,
			},
		},
	})
	assert.NoError(t, err)

	_, err = srv.GetClusterConfig(ctx, &v1.GetClusterConfigRequest{
		ClusterId: c.Id,
	})
	assert.NoError(t, err)

	_, err = srv.DeleteClusterConfig(ctx, &v1.DeleteClusterConfigRequest{
		ClusterId: c.Id,
	})
	assert.NoError(t, err)
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

			srv := New(st, testr.New(t), time.Hour)
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
