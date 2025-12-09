package embedding

type GoogleEmbeddingConfig struct {
	ModelID string `yaml:"model_id"`
	APIKey  string `yaml:"api_key"`
}

func GetDefaultGoogleEmbeddingConfig() *GoogleEmbeddingConfig {
	return &GoogleEmbeddingConfig{
		ModelID: "gemini-2.0-flash",
		APIKey:  "",
	}
}
