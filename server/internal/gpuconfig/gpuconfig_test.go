package gpuconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestCreateTimeSlicingDevicePluginConfig(t *testing.T) {
	c := CreateTimeSlicingDevicePluginConfig(2)
	b, err := yaml.Marshal(c)
	assert.NoError(t, err)
	want := `version: v1
flags:
  commandlineflags:
    migStrategy: none
sharing:
  timeSlicing:
    failRequestsGreaterThanOne: true
    resources:
    - name: nvidia.com/gpu
      replicas: 2
`
	assert.Equal(t, want, string(b))
}
