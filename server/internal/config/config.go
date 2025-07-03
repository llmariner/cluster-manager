package config

import (
	"fmt"
	"os"
	"time"

	"github.com/llmariner/api-usage/pkg/sender"
	"github.com/llmariner/common/pkg/db"
	"gopkg.in/yaml.v3"
)

const defaultComponentStatusTimeout = 1 * time.Hour

// DebugConfig is the debug configuration.
type DebugConfig struct {
	Standalone bool   `yaml:"standalone"`
	SqlitePath string `yaml:"sqlitePath"`
}

// AuthConfig is the authentication configuration.
type AuthConfig struct {
	Enable                 bool   `yaml:"enable"`
	RBACInternalServerAddr string `yaml:"rbacInternalServerAddr"`
}

// Validate validates the configuration.
func (c *AuthConfig) validate() error {
	if !c.Enable {
		return nil
	}
	if c.RBACInternalServerAddr == "" {
		return fmt.Errorf("rbacInternalServerAddr must be set")
	}
	return nil
}

// DefaultClusterConfig is the default cluster configuration.
type DefaultClusterConfig struct {
	Name            string `yaml:"name"`
	RegistrationKey string `yaml:"registrationKey"`
	TenantID        string `yaml:"tenantId"`
}

func (c *DefaultClusterConfig) validate() error {
	if c.Name == "" {
		// Do nothing.
		return nil
	}

	if c.RegistrationKey == "" {
		return fmt.Errorf("registrationKey must be set")
	}
	if c.TenantID == "" {
		return fmt.Errorf("tenantId must be set")
	}
	return nil
}

// NVIDIAConfig is the configuration for Nvidia components.
type NVIDIAConfig struct {
	DevicePluginConfigMapName      string `yaml:"devicePluginConfigMapName"`
	DevicePluginConfigMapNamespace string `yaml:"devicePluginConfigMapNamespace"`
	DevicePluginConfigName         string `yaml:"devicePluginConfigName"`
}

func (c *NVIDIAConfig) validate() error {
	if c.DevicePluginConfigMapName == "" {
		return fmt.Errorf("devicePluginConfigMapName must be set")
	}
	if c.DevicePluginConfigMapNamespace == "" {
		return fmt.Errorf("devicePluginConfigMapNamespace must be set")
	}
	if c.DevicePluginConfigName == "" {
		return fmt.Errorf("devicePluginConfigName must be set")
	}
	return nil
}

// Config is the configuration.
type Config struct {
	GRPCPort              int `yaml:"grpcPort"`
	HTTPPort              int `yaml:"httpPort"`
	WorkerServiceGRPCPort int `yaml:"workerServiceGrpcPort"`
	InternalGRPCPort      int `yaml:"internalGrpcPort"`

	SessionManagerServerEndpoint string `yaml:"sessionManagerServerEndpoint"`

	Database db.Config `yaml:"database"`

	Debug DebugConfig `yaml:"debug"`

	AuthConfig AuthConfig `yaml:"auth"`

	UsageSender sender.Config `yaml:"usageSender"`

	ComponentStatusTimeout time.Duration `yaml:"componentStatusTimeout"`

	DefaultCluster DefaultClusterConfig `yaml:"defaultCluster"`

	NVIDIA NVIDIAConfig `yaml:"nvidia"`
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.GRPCPort <= 0 {
		return fmt.Errorf("grpcPort must be greater than 0")
	}
	if c.HTTPPort <= 0 {
		return fmt.Errorf("httpPort must be greater than 0")
	}
	if c.WorkerServiceGRPCPort <= 0 {
		return fmt.Errorf("workerServiceGrpcPort must be greater than 0")
	}
	if c.InternalGRPCPort <= 0 {
		return fmt.Errorf("internalGrpcPort must be greater than 0")
	}

	if c.SessionManagerServerEndpoint == "" {
		return fmt.Errorf("session manager server endpoint must be set")
	}

	if c.Debug.Standalone {
		if c.Debug.SqlitePath == "" {
			return fmt.Errorf("sqlite path must be set")
		}
	} else {
		if err := c.Database.Validate(); err != nil {
			return fmt.Errorf("database: %s", err)
		}
	}

	if err := c.AuthConfig.validate(); err != nil {
		return err
	}

	if err := c.DefaultCluster.validate(); err != nil {
		return fmt.Errorf("defaultCluster: %s", err)
	}

	if err := c.NVIDIA.validate(); err != nil {
		return fmt.Errorf("nvidia: %s", err)
	}

	if err := c.UsageSender.Validate(); err != nil {
		return err
	}

	if c.ComponentStatusTimeout == 0 {
		c.ComponentStatusTimeout = defaultComponentStatusTimeout
	} else if c.ComponentStatusTimeout < 0 {
		return fmt.Errorf("component status timeout must be greater than 0")
	}

	return nil
}

// Parse parses the configuration file at the given path, returning a new
// Config struct.
func Parse(path string) (Config, error) {
	var config Config

	b, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("config: read: %s", err)
	}

	if err = yaml.Unmarshal(b, &config); err != nil {
		return config, fmt.Errorf("config: unmarshal: %s", err)
	}
	return config, nil
}
