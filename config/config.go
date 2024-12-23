package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/xolinar/kademlia-go/routing"
	"gopkg.in/yaml.v3"
)

// KademliaConfig holds the settings for configuring the Kademlia network.
type Config struct {
	KSize                 routing.KSize `yaml:"ksize" validate:"required,gt=0"`
	Alpha                 int           `yaml:"alpha" validate:"required,gt=0"`
	RepublishInterval     time.Duration `yaml:"republish_interval" validate:"required"`
	RefreshInterval       time.Duration `yaml:"refresh_interval" validate:"required"`
	ReplicationFactor     int           `yaml:"replication_factor" validate:"required,gt=0"`
	MaxConcurrentRequests int           `yaml:"max_concurrent_requests" validate:"required,gt=0"`
	EnableBackgroundTasks bool          `yaml:"enable_background_tasks"`

	NetworkConfig     NetworkConfig     `yaml:"network" validate:"required"`
	TransportConfig   TransportConfig   `yaml:"transport"`
	PerformanceConfig PerformanceConfig `yaml:"performance"`
	LoggingConfig     LoggingConfig     `yaml:"logging"`
}

// TransportConfig holds the configuration for initializing the transport layer in the Kademlia network.
//
// This structure provides the flexibility to specify the transport type (e.g., "udp", "grpc") and any
// protocol-specific parameters, such as IP address, port, or security certificates. It integrates
// seamlessly with the existing configuration system, ensuring that transport settings are applied consistently.
type TransportConfig struct {
	Type     string `yaml:"type" validate:"required"`    // Transport type ("udp", "grpc", etc.)
	Address  string `yaml:"address" validate:"required"` // IP address to bind the transport
	Port     int    `yaml:"port" validate:"required"`    // Port number for the transport
	CertFile string `yaml:"cert_file"`                   // Path to the TLS certificate (optional, for gRPC)
	KeyFile  string `yaml:"key_file"`                    // Path to the TLS key (optional, for gRPC)
}

// NetworkConfig contains network-specific parameters.
type NetworkConfig struct {
	PingTimeout       time.Duration `yaml:"ping_timeout" validate:"required"`
	FindNodeTimeout   time.Duration `yaml:"find_node_timeout" validate:"required"`
	RequestRetryCount int           `yaml:"request_retry_count" validate:"required,gt=0"`
	MaxNodeFailures   int           `yaml:"max_node_failures" validate:"required,gt=0"`
}

// PerformanceConfig includes optimization parameters.
type PerformanceConfig struct {
	ReplacementQueueSize int  `yaml:"replacement_queue_size" validate:"gte=0"`
	EnableOptimizations  bool `yaml:"enable_optimizations"`
}

// LoggingConfig sets the parameters for logging.
type LoggingConfig struct {
	LogLevel  string `yaml:"log_level" validate:"required"`
	LogOutput string `yaml:"log_output" validate:"required"`
}

// LoadConfig reads a YAML configuration file, validates its contents, and returns a KademliaConfig instance.
// The path argument specifies the location of the configuration file.
func LoadConfig(path string) (*Config, error) {
	// Read the configuration file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal YAML into KademliaConfig
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Validate the loaded configuration
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("configuration validation error: %w", err)
	}

	return &config, nil
}
