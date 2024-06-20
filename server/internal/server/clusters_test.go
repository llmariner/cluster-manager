package server

import (
	"context"
	"testing"

	v1 "github.com/llm-operator/cluster-manager/api/v1"
	"github.com/llm-operator/cluster-manager/server/internal/config"
	"github.com/llm-operator/cluster-manager/server/internal/store"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestClusters(t *testing.T) {
	st, tearDown := store.NewTest(t)
	defer tearDown()

	srv := New(st)
	isrv := NewInternal(st)
	ctx := context.Background()

	c, err := srv.CreateCluster(ctx, &v1.CreateClusterRequest{
		Name: "cluster",
	})
	assert.NoError(t, err)
	assert.Equal(t, "cluster", c.Name)
	assert.NotEmpty(t, c.RegistrationKey)

	c, err = srv.GetCluster(ctx, &v1.GetClusterRequest{
		Id: c.Id,
	})
	assert.NoError(t, err)
	assert.Empty(t, c.RegistrationKey)

	listResp, err := srv.ListClusters(ctx, &v1.ListClustersRequest{})
	assert.NoError(t, err)
	assert.Len(t, listResp.Data, 1)
	assert.Equal(t, c.Id, listResp.Data[0].Id)
	assert.Empty(t, listResp.Data[0].RegistrationKey)

	ilistResp, err := isrv.ListInternalClusters(ctx, &v1.ListInternalClustersRequest{})
	assert.NoError(t, err)
	assert.Len(t, ilistResp.Clusters, 1)
	ic := ilistResp.Clusters[0]
	assert.Equal(t, c.Id, ic.Cluster.Id)
	assert.NotEmpty(t, ic.Cluster.RegistrationKey)
	assert.Equal(t, defaultTenantID, ic.TenantId)

	_, err = srv.DeleteCluster(ctx, &v1.DeleteClusterRequest{
		Id: c.Id,
	})
	assert.NoError(t, err)

	_, err = srv.GetCluster(ctx, &v1.GetClusterRequest{
		Id: c.Id,
	})
	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))

	listResp, err = srv.ListClusters(ctx, &v1.ListClustersRequest{})
	assert.NoError(t, err)
	assert.Len(t, listResp.Data, 0)
}

func TestGetSelfCluster(t *testing.T) {
	st, tearDown := store.NewTest(t)
	defer tearDown()

	wsrv := NewWorkerServiceServer(st)
	ctx := context.Background()

	_, err := st.CreateCluster(store.ClusterSpec{
		Name:      "cluster",
		ClusterID: defaultClusterID,
		TenantID:  defaultTenantID,
	})
	assert.NoError(t, err)

	c, err := wsrv.GetSelfCluster(ctx, &v1.GetSelfClusterRequest{})
	assert.NoError(t, err)
	assert.Equal(t, "cluster", c.Name)
}

func TestCreateDefaultCluster(t *testing.T) {
	st, tearDown := store.NewTest(t)
	defer tearDown()

	srv := New(st)

	c := &config.DefaultClusterConfig{
		Name:            "default-cluster",
		RegistrationKey: "rkey",
		TenantID:        "t0",
	}
	err := srv.CreateDefaultCluster(c)
	assert.NoError(t, err)
}
