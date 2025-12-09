package config

type Config struct {
	LLMConfig       LLMConfig       `yaml:"llm_config"`
	LoggerConfig    LoggerConfig    `yaml:"logger_config"`
	TTSConfig       TTSConfig       `yaml:"tts_config"`
	ASRConfig       ASRConfig       `yaml:"asr_config"`
	CharacterConfig CharacterConfig `yaml:"character_config"`
	AgentConfig     AgentConfig     `yaml:"agent_config"`
	ModelsConfig    ModelsConfig    `yaml:"models_config"`
	EmbeddingConfig EmbeddingConfig `yaml:"embedding_config"`
}

func (c *Config) Validate() error {

	if err := c.LLMConfig.Validate(); err != nil {
		return err
	}
	if err := c.TTSConfig.Validate(); err != nil {
		return err
	}
	if err := c.AgentConfig.Validate(); err != nil {
		return err
	}
	if err := c.ModelsConfig.Validate(); err != nil {
		return err
	}
	return nil
}
