package store

import (
	"gorm.io/gorm"
)

// ClusterComponent represents a component in a cluster.
type ClusterComponent struct {
	gorm.Model

	ClusterID     string `gorm:"uniqueIndex:idx_component_cluster_id_name"`
	Name          string `gorm:"uniqueIndex:idx_component_cluster_id_name"`
	IsHealthy     bool
	StatusMessage string
}

// CreateClusterComponent creates a new cluster component in the store
func CreateClusterComponent(tx *gorm.DB, c *ClusterComponent) error {
	res := tx.Create(c)
	return res.Error
}

// UpdateOrCreateClusterComponent sets the appropriate cluster component fields in the
// database given a pointer to a component.
func (s *S) UpdateOrCreateClusterComponent(c *ClusterComponent) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if _, err := FindClusterComponent(tx, c.ClusterID, c.Name); err != nil {
			if err == gorm.ErrRecordNotFound {
				return CreateClusterComponent(tx, c)
			}
			return err
		}

		res := tx.Model(c).
			Where("cluster_id = ? AND name = ?", c.ClusterID, c.Name).
			Update("is_healthy", c.IsHealthy).
			Update("status_message", c.StatusMessage)
		if err := res.Error; err != nil {
			return err
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

// FindClusterComponents fetches the components for the given cluster.
func (s *S) FindClusterComponents(clusterID string) ([]ClusterComponent, error) {
	var cs []ClusterComponent
	if err := s.db.Where("cluster_id = ?", clusterID).Find(&cs).Error; err != nil {
		return nil, err
	}
	return cs, nil
}

// FindClusterComponent fetches the component for the given cluster and component name.
func FindClusterComponent(tx *gorm.DB, clusterID, name string) (*ClusterComponent, error) {
	var c ClusterComponent
	if err := tx.Where("cluster_id = ? AND name = ?", clusterID, name).Take(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// DeleteClusterComponents deletes cluster components by cluster ID.
func DeleteClusterComponents(tx *gorm.DB, clusterID string) error {
	res := tx.Unscoped().Where("cluster_id = ?", clusterID).Delete(&ClusterComponent{})
	if err := res.Error; err != nil {
		return err
	}
	return nil
}
