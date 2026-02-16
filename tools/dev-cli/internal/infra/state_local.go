package infra

import (
	infraStateStore "github.com/useportcall/portcall/tools/dev-cli/internal/infra/state"
)

type ClusterState = infraStateStore.ClusterState
type infraState = infraStateStore.FileState

func infraStatePath(root string) string {
	return infraStateStore.Path(root)
}

func loadInfraState(root string) (infraState, error) {
	return infraStateStore.Load(root)
}

func saveInfraState(root string, state infraState) error {
	return infraStateStore.Save(root, state)
}

func getInfraClusterState(cluster string) (ClusterState, bool) {
	if err := ensureRootDir(); err != nil {
		return ClusterState{}, false
	}
	return LoadClusterState(rootDir, cluster)
}

func LoadClusterState(root, cluster string) (ClusterState, bool) {
	state, err := loadInfraState(root)
	if err != nil {
		return ClusterState{}, false
	}
	cfg, ok := state.Clusters[cluster]
	return cfg, ok
}

func setInfraClusterState(cluster string, cfg ClusterState) error {
	if err := ensureRootDir(); err != nil {
		return err
	}
	state, err := loadInfraState(rootDir)
	if err != nil {
		return err
	}
	state.Clusters[cluster] = cfg
	return saveInfraState(rootDir, state)
}
