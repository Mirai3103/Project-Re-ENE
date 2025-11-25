package services

import (
	"log/slog"

	"github.com/Mirai3103/Project-Re-ENE/config"
)

type ConfigService struct {
	cfg    *config.Config
	logger *slog.Logger
}

func NewConfigService(cfg *config.Config, logger *slog.Logger) *ConfigService {
	return &ConfigService{cfg: cfg, logger: logger}
}

func (h *ConfigService) GetConfig() *config.Config {
	return h.cfg
}

func (h *ConfigService) PatchConfig(cfg *config.Config) {
	log := h.logger
	if err := config.MergeConfig(h.cfg, cfg); err != nil {
		log.Error("merge config", "error", err)
		return
	}
}
