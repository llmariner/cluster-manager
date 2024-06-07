package store

import (
	"gorm.io/gorm"
)

// Cluster represents a cluster.
type Cluster struct {
	gorm.Model

	ClusterID string `gorm:"uniqueIndex"`

	TenantID string `gorm:"uniqueIndex:idx_cluster_tenant_id_name"`

	Name string `gorm:"uniqueIndex:idx_cluster_tenant_id_name"`

	RegistrationKey string
}

// ClusterSpec is a spec of the cluster
type ClusterSpec struct {
	ClusterID       string
	TenantID        string
	Name            string
	RegistrationKey string
}

// CreateCluster creates a cluster.
func (s *S) CreateCluster(spec ClusterSpec) (*Cluster, error) {
	c := &Cluster{
		ClusterID:       spec.ClusterID,
		TenantID:        spec.TenantID,
		Name:            spec.Name,
		RegistrationKey: spec.RegistrationKey,
	}
	if err := s.db.Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

// GetCluster returns a cluster by cluster ID and tenant ID.
func (s *S) GetCluster(clusterID, tenantID string) (*Cluster, error) {
	var c Cluster
	if err := s.db.Where("cluster_id = ? AND tenant_id = ?", clusterID, tenantID).Take(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// ListClustersByTenantID lists clusters.
func (s *S) ListClustersByTenantID(tenantID string) ([]*Cluster, error) {
	var cs []*Cluster
	if err := s.db.Where("tenant_id = ?", tenantID).Order("id DESC").Find(&cs).Error; err != nil {
		return nil, err
	}
	return cs, nil
}

// ListClusters lists clusters.
func (s *S) ListClusters() ([]*Cluster, error) {
	var cs []*Cluster
	if err := s.db.Order("id DESC").Find(&cs).Error; err != nil {
		return nil, err
	}
	return cs, nil
}

// DeleteCluster deletes a cluster by cluster ID and tenant ID.
func (s *S) DeleteCluster(clusterID, tenantID string) error {
	res := s.db.Unscoped().Where("cluster_id = ? AND tenant_id = ?", clusterID, tenantID).Delete(&Cluster{})
	if err := res.Error; err != nil {
		return err
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
