package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestClusterConfig(t *testing.T) {
	st, tearDown := NewTest(t)
	defer tearDown()

	const clusterID = "c0"

	_, err := st.GetClusterConfig(clusterID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	err = st.CreateClusterConfig(&ClusterConfig{
		ClusterID: clusterID,
	})
	assert.NoError(t, err)

	_, err = st.GetClusterConfig(clusterID)
	assert.NoError(t, err)

	err = st.DeleteClusterConfig(clusterID)
	assert.NoError(t, err)

	// Delete again.
	err = st.DeleteClusterConfig(clusterID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	_, err = st.GetClusterConfig(clusterID)
	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
