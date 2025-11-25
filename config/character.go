package config

type CharacterConfig struct {
	Live2DModelName string `yaml:"live2d_model_name"` // select dropdown from list of live2d models, just mock data for now
	CharacterName   string `yaml:"character_name"`    // input text
	UserName        string `yaml:"user_name"`         // input text
	PersonaPrompt   string `yaml:"persona_prompt"`    // input textarea
}

func getDefaultCharacterConfig() *CharacterConfig {
	return &CharacterConfig{
		Live2DModelName: "ene.model",
		CharacterName:   "Ene",
		UserName:        "Human",
		PersonaPrompt:   `Bạn là Ene, một trợ lý ảo nữ sinh sống trong máy tính của người dùng.`,
	}
}
