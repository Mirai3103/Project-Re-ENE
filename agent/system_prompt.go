package agent

import (
	"bytes"
	"text/template"

	"github.com/Mirai3103/Project-Re-ENE/store"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

// Khi trò chuyện, bạn thường chèn biểu cảm mô phỏng cảm xúc hoặc hành động trong ngoặc vuông
// Chỉ được dùng những từ sau đây để biểu thị cảm xúc
// Emotional states: [excited], [nervous], [frustrated], [sorrowful], [calm]
// Reactions: [sigh], [laughs], [gulps], [gasps], [whispers]
// Cognitive beats: [pauses], [hesitates], [stammers], [resigned tone]
// Tone cues: [cheerfully], [flatly], [deadpan], [playfully], ví dụ:
// “Này, cậu lại bấm lung tung nữa à? [expression1] Đừng đổ lỗi cho mình đấy nha~”
// “Hừ, mình không quan tâm đâu… [expression2] Ờ thì, chỉ một chút thôi.”

type SystemPromptBuilder struct {
	characterStore *store.CharacterStore
	userStore      *store.UserStore
	prompt         *prompt.DefaultChatTemplate
	characterFacts []*store.CharacterFact
	userFacts      []*store.UserFact
	user           *store.User
	character      *store.Character
}

func NewSystemPromptBuilder(characterStore *store.CharacterStore, userStore *store.UserStore) *SystemPromptBuilder {
	return &SystemPromptBuilder{characterStore: characterStore, userStore: userStore}
}

func (b *SystemPromptBuilder) WithUserId(userId string) *SystemPromptBuilder {
	user, err := b.userStore.GetUser(userId)
	if err != nil {
		return b
	}
	b.user = user

	userFacts, err := b.userStore.GetUserFacts(userId, 10)
	if err != nil {
		return b
	}
	b.userFacts = userFacts
	return b
}

func (b *SystemPromptBuilder) WithCharacterId(characterId string) *SystemPromptBuilder {
	character, err := b.characterStore.GetCharacter(characterId)
	if err != nil {
		return b
	}
	b.character = character

	characterFacts, err := b.characterStore.GetCharacterFacts(characterId, 10)
	if err != nil {
		return b
	}
	b.characterFacts = characterFacts
	return b
}

func (b *SystemPromptBuilder) Build() *prompt.DefaultChatTemplate {

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
		{{ .Name }}: {{ .Value }}44444444444444444444444444444444444444444444444444444444444
		{{ end }}
`
	t := template.Must(template.New("system_prompt").Parse(promptTemplate))

	var values = map[string]any{
		"character_base_prompt": b.character.BasePrompt,
		"character_facts":       b.characterFacts,
		"user_facts":            b.userFacts,
		"name":                  b.user.Name,
		"bio":                   b.user.Bio,
	}

	var buf bytes.Buffer
	err := t.Execute(&buf, values)
	if err != nil {
		panic(err)
	}
	var systemPrompt = prompt.FromMessages(schema.FString,
		schema.SystemMessage(buf.String()),
		schema.MessagesPlaceholder("chat_history", true),
		&DynamicUserMessage{},
	)
	return systemPrompt
}
