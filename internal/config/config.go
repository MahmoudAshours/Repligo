package config

import (
	"encoding/json"
	"os"
	"repligo/internal/network"
)

type NodeConfig struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type ClusterConfig struct {
	Master network.Network   `json:"master"`
	Slaves []network.Network `json:"slaves"`
}

type Config struct {
	Cluster ClusterConfig `json:"cluster"`
}

func Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
