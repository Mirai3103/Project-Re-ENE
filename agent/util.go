package agent

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/cloudwego/eino/schema"
)

type DynamicUserMessage struct {
}

const (
	USER_INPUT_KEY = "user_input"
	USER_IMAGE_KEY = "user_image"
)

func (d *DynamicUserMessage) Format(ctx context.Context, vs map[string]any, formatType schema.FormatType) ([]*schema.Message, error) {
	fmt.Println("vs: ", vs)
	msg := &schema.Message{
		Role:    schema.User,
		Content: vs[USER_INPUT_KEY].(string),
		UserInputMultiContent: []schema.MessageInputPart{
			{
				Type: schema.ChatMessagePartTypeText,
				Text: vs[USER_INPUT_KEY].(string),
			},
		},
	}

	// Check nếu có image trong params

	if image, ok := vs[USER_IMAGE_KEY].([]byte); ok && image != nil {
		base64Data := base64.StdEncoding.EncodeToString(image)
		base64DataURL := "data:image/png;base64," + base64Data
		msg.UserInputMultiContent = append(msg.UserInputMultiContent, schema.MessageInputPart{
			Type: schema.ChatMessagePartTypeImageURL,
			Image: &schema.MessageInputImage{
				MessagePartCommon: schema.MessagePartCommon{
					URL:        &base64DataURL,
					MIMEType:   "image/png",
					Base64Data: &base64Data,
				},
			},
		})
	}
	fmt.Println("msg: ")
	fmt.Println(msg)
	return []*schema.Message{msg}, nil
}
