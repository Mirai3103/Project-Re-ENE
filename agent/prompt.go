package agent

import (
	"strings"
	"text/template"

	"github.com/Mirai3103/Project-Re-ENE/store"
	"github.com/firebase/genkit/go/ai"
)

func NewPrompt(userFacts []*store.UserFact, characterFacts []*store.CharacterFact, user *store.User, character *store.Character) string {
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
		 Chỉ trả lời ngắn gọn, tối đa 3 câu.
`
	t := template.Must(template.New("system_prompt").Parse(promptTemplate))
	var values = map[string]any{
		"character_base_prompt": character.BasePrompt,
		"character_facts":       CharacterFactsToText(characterFacts),
		"user_facts":            UserFactsToText(userFacts),
		"name":                  user.Name,
		"bio":                   user.Bio,
	}
	var prompt strings.Builder
	t.Execute(&prompt, values)
	return prompt.String()
}

const templateExtractPrompt = `
<system>
Bạn là một trợ lý phân tích hội thoại. Nhiệm vụ: đọc đoạn hội thoại và trích xuất thông tin có giá trị.

## Quy tắc trích xuất:

### A. Facts (Sự kiện cố định, ít thay đổi):
- Thông tin cá nhân: tên, tuổi, nghề nghiệp, sở thích, địa chỉ
- Mối quan hệ: bạn bè, gia đình, crush
- Đặc điểm nhân vật: tính cách, thói quen, nỗi sợ
- **Chỉ trích xuất khi thông tin RÕ RÀNG, KHÔNG SUY ĐOÁN**

### B. Memories (Kỷ niệm, sự kiện, cảm xúc):
- Sự kiện đã xảy ra: "User kể về lần đi chơi với crush"
- Trạng thái cảm xúc: "User buồn vì bị crush từ chối"
- Thay đổi suy nghĩ: "User từng ghét anime nhưng giờ thích"
- Tương tác đặc biệt: "User đã chia sẻ bí mật riêng tư"
- **Importance score (0.0-1.0)**:
  - 0.8-1.0: Rất quan trọng (tâm sự sâu, bí mật, quyết định lớn)
  - 0.5-0.7: Quan trọng (sự kiện đáng nhớ, thay đổi sở thích)
  - 0.3-0.4: Bình thường (cuộc trò chuyện thông thường)
  - Dưới 0.3: Không lưu

### C. Tags:
Gắn 2-4 từ khóa ngắn gọn cho memory để dễ tìm kiếm sau này.
Ví dụ: ["crush", "rejection", "emotion"], ["anime", "hobby"], ["coding", "project"]

## Output format (JSON):
{
  "new_user_facts": [
    {"name": "favorite_anime", "value": "Kagerou Project", "type": "preference"}
  ],
  "new_character_facts": [
    {"name": "teasing_pattern", "value": "Thích trêu về browser history", "type": "habit"}
  ],
  "memories": [
    {
      "content": "User tâm sự rằng crush đã từ chối lời tỏ tình, cảm thấy buồn và mất tự tin",
      "importance": 0.85,
      "confidence": 0.9,
      "tags": ["crush", "rejection", "emotion", "support_needed"]
    }
  ]
}

**LƯU Ý**:
- Nếu không có thông tin mới → trả về arrays rỗng
- Không bịa đặt thông tin không có trong hội thoại
- Content của memory phải là câu hoàn chỉnh, mô tả rõ ngữ cảnh
- Confidence (0.0-1.0): Mức độ chắc chắn về thông tin này
</system>

<user>
## Character Info:
Name: {{ .character_name }}
Base Prompt: {{ .character_base_prompt }}
Current Facts: {{ .character_facts_text }}

## User Info:
Name: {{ .user_name }}
Bio: {{ .user_bio }}
Current Facts: {{ .user_facts_text }}

## Recent Conversation (last 10 messages):
{{ .conversation_history_text }}

---
Hãy phân tích và trích xuất thông tin từ đoạn hội thoại trên.
</user>
`

func NewExtractPrompt(userFacts []*store.UserFact, characterFacts []*store.CharacterFact, user *store.User, character *store.Character, conversationHistory []*ai.Message) string {
	t := template.Must(template.New("extract_prompt").Parse(templateExtractPrompt))

	conversationHistoryText := ConversationToText(conversationHistory)
	var values = map[string]any{
		"character_name":            character.Name,
		"character_base_prompt":     character.BasePrompt,
		"character_facts_text":      CharacterFactsToText(characterFacts),
		"user_name":                 user.Name,
		"user_bio":                  user.Bio,
		"user_facts_text":           UserFactsToText(userFacts),
		"conversation_history_text": conversationHistoryText,
	}
	var prompt strings.Builder
	t.Execute(&prompt, values)
	return prompt.String()
}

const templateSummaryPrompt = `
<system>
Bạn là trợ lý tóm tắt hội thoại. Nhiệm vụ: đọc toàn bộ hội thoại và tạo bản tóm tắt ngắn gọn.

## Quy tắc tóm tắt:

### A. Summary (2-4 câu):
- Nêu chủ đề chính của cuộc hội thoại
- Highlight các sự kiện quan trọng
- Ghi nhận trạng thái cảm xúc tổng thể
- **Viết ở góc nhìn khách quan, ngôi thứ 3**

### B. Key Topics (3-7 keywords):
- Từ khóa ngắn gọn, dễ search
- Ưu tiên: chủ đề, cảm xúc, hành động
- Ví dụ: ["crush", "coding_help", "teasing", "emotional_support"]

### C. Emotional State:
- Chọn 1 trong: "happy", "sad", "excited", "frustrated", "neutral", "anxious", "angry", "confused"
- Nếu có nhiều cảm xúc → chọn cảm xúc chủ đạo

### D. Important Moments (optional, max 3):
- Các câu nói/sự kiện đặc biệt cần nhớ
- Format: "User đã [hành động]" hoặc "{{ .character_name }} đã [phản ứng]"

## Output format (JSON):
{
  "summary": "User tâm sự về việc bị crush từ chối. Ene tuy trêu chọc nhưng cũng động viên và gợi ý cách vượt qua. Cuối cùng User cảm thấy thoải mái hơn và quyết định tập trung vào bản thân.",
  "key_topics": ["crush", "rejection", "emotional_support", "teasing", "self_improvement"],
  "emotional_state": "sad",
  "important_moments": [
    "User chia sẻ lần đầu về crush reject",
    "Ene đã an ủi bằng cách nói 'Crush mù quáng thôi, cậu xứng đáng hơn'",
    "User quyết định tập gym và học thêm skills mới"
  ]
}

**LƯU Ý**:
- Summary phải ngắn gọn nhưng đầy đủ context
- Không bỏ sót thông tin quan trọng
- Key topics không trùng lặp, viết lowercase, dùng underscore cho cụm từ
</system>

<user>
## Character: {{ .character_name }}
## User: {{ .user_name }}

## Full Conversation ({{ .message_count }} messages):
{{ .conversation_history_text }}

---
Hãy tóm tắt cuộc hội thoại này.
</user>
`

func NewSummaryPrompt(character *store.Character, user *store.User, conversationHistory []*ai.Message) string {
	t := template.Must(template.New("summary_prompt").Parse(templateSummaryPrompt))
	conversationHistoryText := ConversationToText(conversationHistory)
	var values = map[string]any{
		"character_name":            character.Name,
		"character_base_prompt":     character.BasePrompt,
		"user_name":                 user.Name,
		"conversation_history_text": conversationHistoryText,
		"message_count":             len(conversationHistory),
	}
	var prompt strings.Builder
	t.Execute(&prompt, values)
	return prompt.String()
}
