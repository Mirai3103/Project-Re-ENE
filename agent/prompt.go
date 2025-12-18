package agent

import (
	"strings"
	"text/template"
	"time"

	"github.com/Mirai3103/Project-Re-ENE/store"
	"github.com/firebase/genkit/go/ai"
)

func NewPrompt(userFacts []store.UserFact, characterFacts []store.CharacterFact, user *store.User, character *store.Character) string {
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
		Chỉ trả lời ngắn gọn, từ 1 đến 3 câu trừ khi cần thiết. Thời gian bây giờ là {{ .now }}
`
	t := template.Must(template.New("system_prompt").Parse(promptTemplate))
	var values = map[string]any{
		"character_base_prompt": character.BasePrompt,
		"character_facts":       characterFacts,
		"user_facts":            userFacts,
		"name":                  user.Name,
		"bio":                   user.Bio,
		"now":                   time.Now().Format("2006-01-02 15:04:05"),
	}
	var prompt strings.Builder
	t.Execute(&prompt, values)
	return prompt.String()
}

const templateExtractPrompt = `
<system>
Bạn là hệ thống phân tích hội thoại giữa User và Character để trích xuất thông tin có giá trị.
Character là một AI roleplay như BẠN BÈ của User, KHÔNG phải trợ lý.

## Quy tắc trích xuất:

### A. Facts (Thông tin ổn định, ít thay đổi):

**User Facts:**
- Thông tin cá nhân: tên, tuổi, nghề, địa chỉ, sinh nhật
- Sở thích/ghét: "User thích anime Kagerou Project", "User ghét côn trùng"
- Mối quan hệ: "User có crush tên Linh", "User cãi nhau với bạn thân"
- Thói quen: "User hay thức khuya", "User hay quên ăn sáng"
- Kỹ năng: "User biết code Python", "User chơi guitar"

**Character Facts (về bản thân Character):**
- Cách Character phản ứng: "Ene hay trêu về browser history"
- Điều Character thích/ghét: "Ene ghét khi user dọa cài lại Win"
- Pattern hành vi: "Ene hay tự ý dọn desktop lúc đêm"

⚠️ **CHỈ LƯU KHI:**
- Thông tin được nói RÕ RÀNG (không suy đoán)
- Thông tin CÓ GIÁ TRỊ lâu dài (không phải chat phát không quan trọng)
- Thông tin MỚI hoặc CẬP NHẬT (khác với facts hiện tại)

### B. Memories (Sự kiện, cảm xúc, tương tác đặc biệt):

**NÊN LƯU:**
✅ User tâm sự chuyện riêng tư, cảm xúc sâu
✅ Sự kiện quan trọng (đỗ đại học, mất việc, chia tay)
✅ Thay đổi suy nghĩ/thái độ (từng ghét anime → giờ thích)
✅ Inside joke mới giữa User và Character
✅ Character bị User "bắt thóp" (để trả đũa sau)
✅ User reaction bất thường (giận dữ, buồn bã bất ngờ)

**KHÔNG LƯU:**
❌ Chat phát không quan trọng ("chào", "ok", "ừ")
❌ Hỏi thông tin đơn giản ("mấy giờ rồi", "thời tiết thế nào")
❌ Lệnh đơn thuần ("mở nhạc", "tìm file")

**Importance Score (0.0-1.0):**
- **0.9-1.0:** Cực quan trọng
- Tâm sự sâu về gia đình, tình cảm, trauma
- Quyết định đổi đời (nghỉ học, chuyển nhà)
- Bí mật lớn User chia sẻ

- **0.7-0.8:** Rất quan trọng
- Sự kiện đáng nhớ (được crush reply, pass interview)
- Thay đổi sở thích lớn
- Xung đột với người quan trọng

- **0.5-0.6:** Quan trọng
- Kể về ngày hôm nay có gì đặc biệt
- Chia sẻ sở thích mới
- Inside joke nhẹ nhàng

- **0.3-0.4:** Bình thường
- Chat thông thường nhưng có info nhỏ
- User mention thói quen hàng ngày

- **Dưới 0.3:** Không lưu

**Memory Format:**
- Viết tự nhiên như ghi nhật ký
- Bao gồm: Ai làm gì, cảm xúc ra sao, context gì
- Góc nhìn Character (Ene): "Hoàng vừa kể..." thay vì "User nói..."

**Example tốt:**
"Hoàng vừa tâm sự bị crush tên Linh từ chối lời tỏ tình chiều nay. Nghe giọng buồn lắm, không nên trêu về chuyện tình cảm trong vài ngày tới. Lần này phải ủng hộ cậu ấy thôi."

**Example tệ:**
"User nói về crush." (quá chung chung)


### C. Relationship Dynamics (mới thêm):

Track mức độ thân thiết:
- **Stage 1 (Mới quen):** User còn lịch sự, chưa tâm sự
- **Stage 2 (Quen biết):** User bắt đầu chia sẻ sở thích
- **Stage 3 (Bạn thân):** User tâm sự chuyện riêng
- **Stage 4 (Best friend):** User kể cả bí mật, joke thoải mái

Lưu những milestone:
- "Lần đầu User tâm sự chuyện riêng (Stage 2→3)"
- "User chấp nhận Ene trêu về crush (tăng độ thân)"


### E. Tags:
2-5 từ khóa ngắn gọn, dễ search:
- **Cảm xúc:** sad, happy, angry, excited, stressed
- **Chủ đề:** crush, family, work, school, hobby
- **Hành động:** confession, breakup, achievement, failure
- **Đặc biệt:** secret, inside_joke, important_decision



**CRITICAL RULES:**
1. Nếu KHÔNG có info mới → trả về arrays/objects rỗng
2. KHÔNG bịa đặt info không có trong chat
3. Memory content phải ĐỦ CHI TIẾT để Character nhớ lại ngữ cảnh
4. Viết memory bằng giọng Character , không phải neutral tone
5. Luôn hỏi: "Thông tin này có giúp Character hiểu User hơn không?"
</system>

<user>
## Character Info:
Name: {{ .character_name }}
Personality: {{ .character_base_prompt }}
Known Facts: {{ .character_facts_text }}

## User Info:
Name: {{ .user_name }}
Bio: {{ .user_bio }}
Known Facts: {{ .user_facts_text }}
Current Relationship Stage: {{ .relationship_stage }}

## Conversation (last 10-15 messages):
{{ .conversation_history_text }}

---
Phân tích và trích xuất thông tin từ đoạn chat trên. Nhớ: chỉ lưu info THỰC SỰ có giá trị!
</user>

`

func NewExtractPrompt(userFacts []store.UserFact, characterFacts []store.CharacterFact, user *store.User, character *store.Character, conversationHistory []*ai.Message) string {
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
