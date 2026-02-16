package initcmd

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type tfOutputEntry struct {
	Value any `json:"value"`
}

type tfOutputs struct {
	RegistryEndpoint                                                   string
	PostgresHost, PostgresPort, PostgresUser, PostgresPassword         string
	PostgresDatabase, KeycloakDatabase                                 string
	RedisHost, RedisPort, RedisPassword                                string
	SpacesRegion, SpacesEndpoint, SpacesAccessKey, SpacesSecretKey     string
	SpacesIconBucket, SpacesQuoteBucket                                string
}

func readTerraformOutputs(plan Plan, deps Deps) (tfOutputs, error) {
	out, err := deps.RunCmdOut("terraform", "-chdir="+plan.StackDir, "output", "-json")
	if err != nil {
		return tfOutputs{}, fmt.Errorf("terraform output failed: %w", err)
	}
	raw := map[string]tfOutputEntry{}
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		return tfOutputs{}, fmt.Errorf("parse terraform outputs: %w", err)
	}
	return tfOutputs{
		RegistryEndpoint: asString(raw, "registry_endpoint"),
		PostgresHost:     asString(raw, "postgres_host"),
		PostgresPort:     asString(raw, "postgres_port"),
		PostgresUser:     asString(raw, "postgres_user"),
		PostgresPassword: asString(raw, "postgres_password"),
		PostgresDatabase: asString(raw, "postgres_database"),
		KeycloakDatabase: asString(raw, "keycloak_database"),
		RedisHost:        asString(raw, "redis_host"),
		RedisPort:        asString(raw, "redis_port"),
		RedisPassword:    asString(raw, "redis_password"),
		SpacesRegion:     asString(raw, "spaces_region"),
		SpacesEndpoint:   asString(raw, "spaces_endpoint"),
		SpacesAccessKey:  asString(raw, "spaces_access_key"),
		SpacesSecretKey:  asString(raw, "spaces_secret_key"),
		SpacesIconBucket: asString(raw, "spaces_icon_bucket"),
		SpacesQuoteBucket: asString(raw, "spaces_quote_bucket"),
	}, nil
}

func asString(raw map[string]tfOutputEntry, key string) string {
	v, ok := raw[key]
	if !ok || v.Value == nil {
		return ""
	}
	switch t := v.Value.(type) {
	case string:
		return t
	case float64:
		return strconv.Itoa(int(t))
	default:
		return fmt.Sprintf("%v", t)
	}
}

func validateManagedOutputs(plan Plan, opts Options, outs tfOutputs) error {
	missing := []string{}
	if !opts.SkipPostgres {
		for k, v := range map[string]string{
			"postgres_host": outs.PostgresHost, "postgres_user": outs.PostgresUser,
			"postgres_password": outs.PostgresPassword, "postgres_port": outs.PostgresPort,
		} {
			if v == "" {
				missing = append(missing, k)
			}
		}
	}
	if !opts.SkipRedis {
		for k, v := range map[string]string{
			"redis_host": outs.RedisHost, "redis_password": outs.RedisPassword, "redis_port": outs.RedisPort,
		} {
			if v == "" {
				missing = append(missing, k)
			}
		}
	}
	if !opts.SkipSpaces {
		for k, v := range map[string]string{
			"spaces_access_key": outs.SpacesAccessKey, "spaces_secret_key": outs.SpacesSecretKey,
		} {
			if v == "" {
				missing = append(missing, k)
			}
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing managed outputs: %v (step=%s)", missing, plan.Step)
	}
	return nil
}
