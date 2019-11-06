package notification

import (
	"context"
	"fmt"
	"html/template"
)

// PushMessageTemplateFile to store template in a file
type PushMessageTemplateFile struct {
	Title         string `json:"title"`
	Message       string `json:"message"`
	MessageDetail string `json:"message_detail"`
	Image         string `json:"image"`
	WebViewLink   string `json:"webview_link"`
}

// PushMessageGoTemplate for storing information about push message template as go template
type PushMessageGoTemplate struct {
	Title         *template.Template
	Message       *template.Template
	MessageDetail *template.Template
	Image         string
	WebViewLink   string
}

// PushMessageTemplate for notification
type PushMessageTemplate struct {
	templates map[string]PushMessageGoTemplate
}

// NewPushMessageTemplate for push notification template
func NewPushMessageTemplate(ctx context.Context, templates map[string]PushMessageTemplateFile) (*PushMessageTemplate, error) {
	goTemplates := make(map[string]PushMessageGoTemplate)

	// TODO: make parsing concurrent with goroutines
	for name, file := range templates {
		var err error
		title := template.New(fmt.Sprintf("%s%s", "title", name))
		title, err = title.Parse(file.Title)
		if err != nil {
			return nil, err
		}

		message := template.New(fmt.Sprintf("%s%s", "message", name))
		message, err = message.Parse(file.Message)
		if err != nil {
			return nil, err
		}

		messageDetail := template.New(fmt.Sprintf("%s%s", "message_detail", name))
		messageDetail, err = messageDetail.Parse(file.MessageDetail)
		if err != nil {
			return nil, err
		}

		goTemplates[name] = PushMessageGoTemplate{
			Title:         title,
			Message:       message,
			MessageDetail: messageDetail,
			Image:         file.Image,
			WebViewLink:   file.WebViewLink,
		}
	}

	t := PushMessageTemplate{
		templates: goTemplates,
	}
	return &t, nil
}

// Execute push message template
func (pmt *PushMessageTemplate) Execute(name string, data interface{}) error {
	_, ok := pmt.templates[name]
	if !ok {
		return fmt.Errorf("push_message_tempalte: template with name %s not found", name)
	}
	return nil
}
