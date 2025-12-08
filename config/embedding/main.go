package embedding

type GoogleGeminiConfig struct {
	ModelID string `yaml:"model_id"`
	APIKey  string `yaml:"api_key"`
}

func GetDefaultGoogleGeminiConfig() *GoogleGeminiConfig {
	return &GoogleGeminiConfig{
		ModelID: "gemini-2.0-flash",
		APIKey:  "",
	}
}
