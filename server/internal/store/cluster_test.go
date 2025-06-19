package store

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCluster(t *testing.T) {
	st, tearDown := NewTest(t)
	defer tearDown()

	const (
		clusterID = "f0"
		tenantID  = "tid0"
	)

	_, err := st.GetCluster(clusterID, tenantID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	_, err = st.CreateCluster(ClusterSpec{
		ClusterID: clusterID,
		TenantID:  tenantID,
		Name:      "name",
	})
	assert.NoError(t, err)

	gotC, err := st.GetCluster(clusterID, tenantID)
	assert.NoError(t, err)
	assert.Equal(t, clusterID, gotC.ClusterID)
	assert.Equal(t, tenantID, gotC.TenantID)

	gotC, err = st.GetClusterByNameAndTenantID("name", tenantID)
	assert.NoError(t, err)
	assert.Equal(t, clusterID, gotC.ClusterID)
	assert.Equal(t, tenantID, gotC.TenantID)

	_, err = st.GetClusterByNameAndTenantID("different_name", tenantID)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	gotCs, err := st.ListClustersByTenantID(tenantID)
	assert.NoError(t, err)
	assert.Len(t, gotCs, 1)

	cs, err := st.FindClusterComponents(clusterID)
	assert.NoError(t, err)
	assert.Len(t, cs, 5)

	_, err = st.CreateCluster(ClusterSpec{
		ClusterID: "f1",
		TenantID:  "tid1",
		Name:      "name",
	})
	assert.NoError(t, err)

	gotCs, err = st.ListClustersByTenantID(tenantID)
	assert.NoError(t, err)
	assert.Len(t, gotCs, 1)

	err = st.DeleteCluster(clusterID, tenantID)
	assert.NoError(t, err)

	gotCs, err = st.ListClustersByTenantID(tenantID)
	assert.NoError(t, err)
	assert.Len(t, gotCs, 0)

	err = st.DeleteCluster(clusterID, tenantID)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

	cs, err = st.FindClusterComponents(clusterID)
	assert.NoError(t, err)
	assert.Len(t, cs, 0)
}
