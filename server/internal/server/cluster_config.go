package server

import (
	"context"
	"errors"

	v1 "github.com/llmariner/cluster-manager/api/v1"
	"github.com/llmariner/cluster-manager/server/internal/store"
	gerrors "github.com/llmariner/common/pkg/gormlib/errors"
	"github.com/llmariner/rbac-manager/pkg/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// CreateClusterConfig creates a cluster config.
func (s *S) CreateClusterConfig(
	ctx context.Context,
	req *v1.CreateClusterConfigRequest,
) (*v1.ClusterConfig, error) {
	userInfo, ok := auth.ExtractUserInfoFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to authorize")
	}

	if err := s.validateClusterID(req.ClusterId, userInfo); err != nil {
		return nil, err
	}

	config := &store.ClusterConfig{
		ClusterID: req.ClusterId,
	}
	if err := s.store.CreateClusterConfig(config); err != nil {
		if gerrors.IsUniqueConstraintViolation(err) {
			return nil, status.Errorf(codes.AlreadyExists, "config for cluster %q already exists", req.ClusterId)
		}
		return nil, status.Errorf(codes.Internal, "create cluster config: %s", err)
	}

	return toConfigProto(config), nil
}

// GetClusterConfig retrieves the cluster config.
func (s *S) GetClusterConfig(
	ctx context.Context,
	req *v1.GetClusterConfigRequest,
) (*v1.ClusterConfig, error) {
	userInfo, ok := auth.ExtractUserInfoFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to authorize")
	}

	if err := s.validateClusterID(req.ClusterId, userInfo); err != nil {
		return nil, err
	}

	config, err := s.store.GetClusterConfig(req.ClusterId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "config for cluster %q not found", req.ClusterId)
		}
		return nil, status.Errorf(codes.Internal, "get cluster config: %s", err)
	}

	return toConfigProto(config), nil
}

// UpdateClusterConfig updates the cluster config.
func (s *S) UpdateClusterConfig(
	ctx context.Context,
	req *v1.UpdateClusterConfigRequest,
) (*v1.ClusterConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "update cluster config is not implemented yet")
}

// DeleteClusterConfig deletes the cluster config.
func (s *S) DeleteClusterConfig(
	ctx context.Context,
	req *v1.DeleteClusterConfigRequest,
) (*emptypb.Empty, error) {
	userInfo, ok := auth.ExtractUserInfoFromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to authorize")
	}

	if err := s.validateClusterID(req.ClusterId, userInfo); err != nil {
		return nil, err
	}

	if err := s.store.DeleteClusterConfig(req.ClusterId); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "config for cluster %q not found", req.ClusterId)
		}
		return nil, status.Errorf(codes.Internal, "delete cluster config: %s", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *S) validateClusterID(clusterID string, userInfo *auth.UserInfo) error {
	if clusterID == "" {
		return status.Errorf(codes.InvalidArgument, "cluster ID is required")
	}

	// Check if the cluster exists for the tenant.
	if _, err := s.store.GetCluster(clusterID, userInfo.TenantID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Errorf(codes.NotFound, "cluster %q not found", clusterID)
		}
		return status.Errorf(codes.Internal, "get cluster: %s", err)
	}

	return nil
}

func toConfigProto(config *store.ClusterConfig) *v1.ClusterConfig {
	return &v1.ClusterConfig{
		// TODO(kenji): Add fields.
	}
}
