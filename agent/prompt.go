package agent

import (
	"strings"
	"text/template"

	"github.com/Mirai3103/Project-Re-ENE/store"
)

func NewPrompt(characterStore *store.CharacterStore, userStore *store.UserStore, conversationStore *store.ConversationStore, input *FlowInput) string {
	var promptTemplate = `
		{{ .character_base_prompt }}
		Đây là những thông tin cá nhân của bạn:
		{{ range .character_facts }}
		{{ .Name }}: {{ .Value }}
		{{ end }}
		Đây là những thông tin cá nhân của người dùng:
		Tên: {{ .name }}
		{{ .bio }}
		{{ range .user_facts }}
		{{ .Name }}: {{ .Value }}
		{{ end }}
`
	t := template.Must(template.New("system_prompt").Parse(promptTemplate))
	userFacts, _ := userStore.GetUserFacts(input.UserID, 10)

	characterFacts, _ := characterStore.GetCharacterFacts(input.CharacterID, 10)
	user, _ := userStore.GetUser(input.UserID)
	character, _ := characterStore.GetCharacter(input.CharacterID)
	var values = map[string]any{
		"character_base_prompt": character.BasePrompt,
		"character_facts":       characterFacts,
		"user_facts":            userFacts,
		"name":                  user.Name,
		"bio":                   user.Bio,
	}
	var prompt strings.Builder
	t.Execute(&prompt, values)
	return prompt.String()
}
