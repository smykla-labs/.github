// Package config provides sync configuration parsing.
// Type definitions are in internal/configtypes for minimal import footprint.
package config

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
	"gopkg.in/yaml.v3"

	"github.com/smykla-labs/.github/internal/configtypes"
)

// ParseSyncConfig parses sync configuration from YAML or JSON.
func ParseSyncConfig(data []byte) (*configtypes.SyncConfig, error) {
	var cfg configtypes.SyncConfig

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		if jsonErr := json.Unmarshal(data, &cfg); jsonErr != nil {
			return nil, errors.Wrap(err, "parsing sync config as YAML or JSON")
		}
	}

	return &cfg, nil
}

// ParseSyncConfigJSON parses sync configuration from JSON string.
func ParseSyncConfigJSON(jsonStr string) (*configtypes.SyncConfig, error) {
	if jsonStr == "" {
		return &configtypes.SyncConfig{}, nil
	}

	var cfg configtypes.SyncConfig
	if err := json.Unmarshal([]byte(jsonStr), &cfg); err != nil {
		return nil, errors.Wrap(err, "parsing sync config JSON")
	}

	return &cfg, nil
}
