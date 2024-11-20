package server

import (
	"context"
	"errors"
	"fmt"

	v1 "github.com/llmariner/cluster-manager/api/v1"
	"github.com/llmariner/cluster-manager/server/internal/config"
	"github.com/llmariner/cluster-manager/server/internal/store"
	gerrors "github.com/llmariner/common/pkg/gormlib/errors"
	"github.com/llmariner/common/pkg/id"
	"github.com/llmariner/rbac-manager/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// CreateCluster creates a cluster.
func (s *S) CreateCluster(
	ctx context.Context,
	req *v1.CreateClusterRequest,
) (*v1.Cluster, error) {
	userInfo, ok := auth.ExtractUserInfoFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract user info from context")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	clusterID, err := id.GenerateID("cluster-", 24)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate cluster id: %s", err)
	}

	rkey, err := id.GenerateID("clusterkey-", 24)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate cluster registration key: %s", err)
	}

	c, err := s.store.CreateCluster(store.ClusterSpec{
		ClusterID:       clusterID,
		TenantID:        userInfo.TenantID,
		Name:            req.Name,
		RegistrationKey: rkey,
	})
	if err != nil {
		if gerrors.IsUniqueConstraintViolation(err) {
			return nil, status.Errorf(codes.AlreadyExists, "cluster %q already exists", req.Name)
		}
		return nil, status.Errorf(codes.Internal, "create cluster: %s", err)
	}

	return toClusterProto(c, true), nil
}

// ListClusters lists clusters.
func (s *S) ListClusters(
	ctx context.Context,
	req *v1.ListClustersRequest,
) (*v1.ListClustersResponse, error) {
	userInfo, ok := auth.ExtractUserInfoFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract user info from context")
	}

	cs, err := s.store.ListClustersByTenantID(userInfo.TenantID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list clusters: %s", err)
	}

	var clusterProtos []*v1.Cluster
	for _, c := range cs {
		cProto := toClusterProto(c, false)
		coms, err := s.store.FindClusterComponents(c.ClusterID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "find cluster components: %s", err)
		}
		cProto.ComponentsStatuses = toComponentStatusesProto(coms)
		clusterProtos = append(clusterProtos, cProto)
	}
	return &v1.ListClustersResponse{
		Object: "list",
		Data:   clusterProtos,
	}, nil
}

// GetCluster gets a cluster.
func (s *S) GetCluster(
	ctx context.Context,
	req *v1.GetClusterRequest,
) (*v1.Cluster, error) {
	userInfo, ok := auth.ExtractUserInfoFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract user info from context")
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	c, err := s.store.GetCluster(req.Id, userInfo.TenantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "cluster %q not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "get cluster: %s", err)
	}
	cProto := toClusterProto(c, false)
	coms, err := s.store.FindClusterComponents(c.ClusterID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "find cluster components: %s", err)
	}
	cProto.ComponentsStatuses = toComponentStatusesProto(coms)
	return cProto, nil
}

// DeleteCluster deletes a cluster.
func (s *S) DeleteCluster(
	ctx context.Context,
	req *v1.DeleteClusterRequest,
) (*v1.DeleteClusterResponse, error) {
	userInfo, ok := auth.ExtractUserInfoFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to extract user info from context")
	}

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.store.DeleteCluster(req.Id, userInfo.TenantID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "cluster %q not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "delete cluster: %s", err)
	}
	return &v1.DeleteClusterResponse{
		Id:      req.Id,
		Object:  "cluster",
		Deleted: true,
	}, nil
}

// CreateDefaultCluster creates a default cluster if it does not exist.
func (s *S) CreateDefaultCluster(c *config.DefaultClusterConfig) error {
	_, err := s.store.GetClusterByNameAndTenantID(c.Name, c.TenantID)
	if err == nil {
		// Do nothing.
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	s.log.Info("Creating default cluster", "name", c.Name)
	clusterID, err := id.GenerateID("cluster-", 24)
	if err != nil {
		return err
	}
	if _, err := s.store.CreateCluster(store.ClusterSpec{
		ClusterID:       clusterID,
		TenantID:        c.TenantID,
		Name:            c.Name,
		RegistrationKey: c.RegistrationKey,
	}); err != nil {
		return err
	}

	return nil
}

// GetSelfCluster gets a cluster where the worker cluster itself belongs.
func (s *WS) GetSelfCluster(
	ctx context.Context,
	req *v1.GetSelfClusterRequest,
) (*v1.Cluster, error) {
	clusterInfo, err := s.extractClusterInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	c, err := s.store.GetCluster(clusterInfo.ClusterID, clusterInfo.TenantID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get cluster: %s", err)
	}
	return toClusterProto(c, false), nil
}

// UpdateComponentStatus updates the status of a component in a cluster.
func (s *WS) UpdateComponentStatus(
	ctx context.Context,
	req *v1.UpdateComponentStatusRequest,
) (*emptypb.Empty, error) {
	clusterInfo, err := s.extractClusterInfoFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if req.Status == nil {
		return nil, status.Error(codes.InvalidArgument, "status is required")
	}

	c := &store.ClusterComponent{
		ClusterID:     clusterInfo.ClusterID,
		Name:          req.Name,
		IsHealthy:     req.Status.IsHealthy,
		StatusMessage: req.Status.Message,
	}

	err = s.store.UpdateOrCreateClusterComponent(c)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update or create cluster component: %s", err)
	}
	return &emptypb.Empty{}, nil
}

// ListInternalClusters lists all clusters with registration keys.
func (s *IS) ListInternalClusters(
	ctx context.Context,
	req *v1.ListInternalClustersRequest,
) (*v1.ListInternalClustersResponse, error) {
	cs, err := s.store.ListClusters()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list clusters: %s", err)
	}

	var clusterProtos []*v1.InternalCluster
	for _, c := range cs {
		clusterProtos = append(clusterProtos, toInternalClusterProto(c))
	}
	return &v1.ListInternalClustersResponse{
		Clusters: clusterProtos,
	}, nil
}

func toInternalClusterProto(c *store.Cluster) *v1.InternalCluster {
	return &v1.InternalCluster{
		Cluster:  toClusterProto(c, true),
		TenantId: c.TenantID,
	}
}

func toClusterProto(c *store.Cluster, withRegistrationKey bool) *v1.Cluster {
	var rkey string
	if withRegistrationKey {
		rkey = c.RegistrationKey
	}

	return &v1.Cluster{
		Id:              c.ClusterID,
		Name:            c.Name,
		RegistrationKey: rkey,
		Object:          "cluster",
	}
}

func toComponentStatusesProto(cs []store.ClusterComponent) map[string]*v1.ComponentStatus {
	m := make(map[string]*v1.ComponentStatus)
	for _, c := range cs {
		m[c.Name] = &v1.ComponentStatus{
			IsHealthy: c.IsHealthy,
			Message:   c.StatusMessage,
		}
	}
	return m
}
