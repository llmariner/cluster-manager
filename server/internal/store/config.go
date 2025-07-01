package store

import "gorm.io/gorm"

// ClusterConfig represents the configuration for a cluster.
type ClusterConfig struct {
	gorm.Model

	ClusterID string `gorm:"uniqueIndex"`

	// Message is the encoded ClusterConfig proto message.
	Message []byte
}

// CreateClusterConfig creates a new ClusterConfig in the database.
func (s *S) CreateClusterConfig(config *ClusterConfig) error {
	if err := s.db.Create(config).Error; err != nil {
		return err
	}
	return nil
}

// GetClusterConfig retrieves a ClusterConfig its ClusterID.
func (s *S) GetClusterConfig(clusterID string) (*ClusterConfig, error) {
	var config ClusterConfig
	if err := s.db.Where("cluster_id = ?", clusterID).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// DeleteClusterConfig deletes a ClusterConfig by its ClusterID.
func (s *S) DeleteClusterConfig(clusterID string) error {
	return deleteClusterConfigInTransaction(s.db, clusterID)
}

func deleteClusterConfigInTransaction(tx *gorm.DB, clusterID string) error {
	res := tx.Unscoped().Where("cluster_id = ?", clusterID).Delete(&ClusterConfig{})
	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
