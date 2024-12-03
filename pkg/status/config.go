package status

import (
	"fmt"
	"time"
)

const (
	defaultInitialDelay = 5 * time.Minute
	defaultInterval     = 1 * time.Hour
)

// Config is the configuration for the sender.
type Config struct {
	// Enable is the flag to enable sending component status.
	Enable bool `yaml:"enable"`
	// Name is the name of the component.
	Name string `yaml:"name"`
	// ClusterManagerServerWorkerServiceAddr is the address of the cluster manager server worker service.
	ClusterManagerServerWorkerServiceAddr string `yaml:"clusterManagerServerWorkerServiceAddr"`
	// InitalDelay is the inital delay at which the sender sends the first status message to the server.
	InitalDelay time.Duration `yaml:"initialDelay"`
	// Interval is the interval at which the sender sends status message to the server.
	Interval time.Duration `yaml:"interval"`
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if !c.Enable {
		return nil
	}
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.ClusterManagerServerWorkerServiceAddr == "" {
		return fmt.Errorf("cluster manager server address is required")
	}
	if c.InitalDelay == 0 {
		c.InitalDelay = defaultInitialDelay
	} else if c.InitalDelay < 0 {
		return fmt.Errorf("inital delay must be greater than 0")
	}
	if c.Interval == 0 {
		c.Interval = defaultInterval
	} else if c.Interval < 0 {
		return fmt.Errorf("interval must be greater than 0")
	}
	return nil
}
