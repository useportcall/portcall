package initcmd

import "fmt"

type bootstrapSecrets struct {
	Portcall map[string]string
	Spaces   map[string]string
}

func finalizeInfra(plan Plan, opts Options, deps Deps) error {
	ctx, err := configureClusterContext(plan, deps)
	if err != nil {
		return err
	}
	outs, err := readTerraformOutputs(plan, deps)
	if err != nil {
		return err
	}
	if plan.Step == "core" {
		if err := bootstrapCluster(plan, opts, bootstrapSecrets{}, deps); err != nil {
			return err
		}
		cfg, _ := deps.GetClusterState(plan.Alias)
		return deps.SetClusterState(plan.Alias, ClusterState{
			Context: ctx, Registry: firstNonEmpty(outs.RegistryEndpoint, plan.Registry), Values: cfg.Values,
			Mode: plan.Mode, Cluster: plan.ClusterName, Namespace: plan.Namespace, Provider: plan.Provider, InitAction: plan.Action,
		})
	}
	if err := validateManagedOutputs(plan, opts, outs); err != nil {
		return err
	}
	vals, secrets, err := writeMicroConfig(plan, outs, deps)
	if err != nil {
		return err
	}
	if err := bootstrapCluster(plan, opts, secrets, deps); err != nil {
		return err
	}
	return deps.SetClusterState(plan.Alias, ClusterState{
		Context: ctx, Registry: outs.RegistryEndpoint, Values: vals,
		Mode: plan.Mode, Cluster: plan.ClusterName, Namespace: plan.Namespace, Provider: plan.Provider, InitAction: plan.Action,
	})
}

func configureClusterContext(plan Plan, deps Deps) (string, error) {
	if err := deps.RunCmd("doctl", "kubernetes", "cluster", "kubeconfig", "save", plan.ClusterName); err != nil {
		return "", fmt.Errorf("save kubeconfig for cluster %s: %w", plan.ClusterName, err)
	}
	ctx, err := deps.RunCmdOut("kubectl", "config", "current-context")
	if err != nil {
		return "", fmt.Errorf("read kubectl current context: %w", err)
	}
	return ctx, nil
}

func writeMicroConfig(plan Plan, outs tfOutputs, deps Deps) (string, bootstrapSecrets, error) {
	aesKey, err := ensureStableSecret(plan.Namespace, "portcall-secrets", "AES_ENCRYPTION_KEY", 32, deps)
	if err != nil {
		return "", bootstrapSecrets{}, err
	}
	kcPass, err := ensureStableSecret(plan.Namespace, "portcall-secrets", "KC_BOOTSTRAP_ADMIN_PASSWORD", 24, deps)
	if err != nil {
		return "", bootstrapSecrets{}, err
	}
	grafanaPass, err := ensureStableSecret(plan.Namespace, "portcall-secrets", "GRAFANA_ADMIN_PASSWORD", 24, deps)
	if err != nil {
		return "", bootstrapSecrets{}, err
	}
	kcUser := firstNonEmpty(getSecretValue(plan.Namespace, "portcall-secrets", "KC_BOOTSTRAP_ADMIN_USERNAME", deps), "admin")
	registry := firstNonEmpty(outs.RegistryEndpoint, plan.Registry)
	pgDB := firstNonEmpty(outs.PostgresDatabase, "main_portcall_db")
	kcDB := firstNonEmpty(outs.KeycloakDatabase, "keycloak")
	path, err := writeMicroValuesFile(plan.Alias, microValues{
		Registry: registry, Domain: plan.Domain,
		PostgresHost: outs.PostgresHost, PostgresPort: outs.PostgresPort, PostgresUser: outs.PostgresUser,
		PostgresDB: pgDB, KeycloakDB: kcDB, RedisHost: outs.RedisHost, RedisPort: outs.RedisPort,
		SpacesRegion: outs.SpacesRegion, SpacesEndpoint: outs.SpacesEndpoint,
		SpacesBucket: firstNonEmpty(outs.SpacesQuoteBucket, plan.SpacesPrefix+"-quote-signatures"),
		AllowedIPs: plan.AllowedIPs,
	}, deps)
	if err != nil {
		return "", bootstrapSecrets{}, err
	}
	pgURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=require",
		outs.PostgresUser, outs.PostgresPassword, outs.PostgresHost, outs.PostgresPort, pgDB)
	portcall := map[string]string{
		"AES_ENCRYPTION_KEY": aesKey, "KC_BOOTSTRAP_ADMIN_USERNAME": kcUser,
		"KC_BOOTSTRAP_ADMIN_PASSWORD": kcPass, "GRAFANA_ADMIN_PASSWORD": grafanaPass,
	}
	putIfValue(portcall, "POSTGRES_USER", outs.PostgresUser)
	putIfValue(portcall, "POSTGRES_PASSWORD", outs.PostgresPassword)
	putIfValue(portcall, "DATABASE_URL", pgURL)
	putIfValue(portcall, "KC_DB_USERNAME", outs.PostgresUser)
	putIfValue(portcall, "KC_DB_PASSWORD", outs.PostgresPassword)
	putIfValue(portcall, "REDIS_PASSWORD", outs.RedisPassword)
	putIfValue(portcall, "S3_ACCESS_KEY_ID", outs.SpacesAccessKey)
	putIfValue(portcall, "S3_SECRET_ACCESS_KEY", outs.SpacesSecretKey)
	spaces := map[string]string{}
	putIfValue(spaces, "AWS_ACCESS_KEY_ID", outs.SpacesAccessKey)
	putIfValue(spaces, "AWS_SECRET_ACCESS_KEY", outs.SpacesSecretKey)
	return path, bootstrapSecrets{Portcall: portcall, Spaces: spaces}, nil
}

func firstNonEmpty(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}

func putIfValue(target map[string]string, key, value string) {
	if value != "" {
		target[key] = value
	}
}
