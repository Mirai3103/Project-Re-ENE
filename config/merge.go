package config

import "dario.cat/mergo"

func MergeConfig(dst *Config, patch *Config) error {
	return mergo.Merge(dst, patch, mergo.WithOverride)
}
