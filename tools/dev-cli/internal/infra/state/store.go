package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type ClusterState struct {
	Context    string `json:"context,omitempty"`
	Registry   string `json:"registry,omitempty"`
	Values     string `json:"values,omitempty"`
	Mode       string `json:"mode,omitempty"`
	Cluster    string `json:"cluster,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Provider   string `json:"provider,omitempty"`
	InitAction string `json:"initAction,omitempty"`
}

type FileState struct {
	Clusters map[string]ClusterState `json:"clusters,omitempty"`
}

func Path(root string) string {
	return filepath.Join(root, ".dev-cli.infra.json")
}

func Load(root string) (FileState, error) {
	var out FileState
	data, err := os.ReadFile(Path(root))
	if errors.Is(err, os.ErrNotExist) {
		out.Clusters = map[string]ClusterState{}
		return out, nil
	}
	if err != nil {
		return out, fmt.Errorf("read infra state: %w", err)
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return out, fmt.Errorf("parse infra state: %w", err)
	}
	if out.Clusters == nil {
		out.Clusters = map[string]ClusterState{}
	}
	return out, nil
}

func Save(root string, state FileState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encode infra state: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(Path(root), data, 0o644); err != nil {
		return fmt.Errorf("write infra state: %w", err)
	}
	return nil
}
