package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComponent(t *testing.T) {
	st, tearDown := NewTest(t)
	defer tearDown()

	const clusterID = "f0"

	err := CreateClusterComponent(st.DB(), &ClusterComponent{
		ClusterID:     clusterID,
		Name:          "name1",
		IsHealthy:     true,
		StatusMessage: "status healthy",
	})
	assert.NoError(t, err)

	err = st.UpdateOrCreateClusterComponent(&ClusterComponent{
		ClusterID:     clusterID,
		Name:          "name1",
		IsHealthy:     false,
		StatusMessage: "status unhealthy",
	})
	assert.NoError(t, err)

	err = st.UpdateOrCreateClusterComponent(&ClusterComponent{
		ClusterID:     clusterID,
		Name:          "name2",
		IsHealthy:     true,
		StatusMessage: "status healthy",
	})
	assert.NoError(t, err)

	got, err := st.FindClusterComponents(clusterID)
	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, clusterID, got[0].ClusterID)
	assert.Equal(t, clusterID, got[1].ClusterID)
	for _, c := range got {
		switch c.Name {
		case "name1":
			assert.False(t, c.IsHealthy)
			assert.Equal(t, "status unhealthy", c.StatusMessage)
		case "name2":
			assert.True(t, c.IsHealthy)
			assert.Equal(t, "status healthy", c.StatusMessage)
		default:
			t.Fatalf("unexpected component: %v", c)
		}
	}
}
