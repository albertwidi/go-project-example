package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	entity "github.com/albertwidi/go-project-example/internal/entity/notification"
)
// notificationTemplate struct
type notificationTemplate struct {
	pushMessageTemplate map[string]*template.Template
	smsTemplate         map[string]*template.Template
	emailTemplate       map[string]*template.Template
}

// newTemplateGen for notification template
func newTemplateGen(basePath string) (*notificationTemplate, error) {
	templatesList := []string{"sms", "pushmessage"}
	ntm := notificationTemplate{
		pushMessageTemplate: make(map[string]*template.Template),
		smsTemplate:         make(map[string]*template.Template),
		emailTemplate:       make(map[string]*template.Template),
	}

	// load the template
	for _, l := range templatesList {
		var err error
		templatePath := path.Join(basePath, l)

		switch l {
		case "sms":
			err = filepath.Walk(templatePath, ntm.loadSMS)
		}

		if err != nil {
			return nil, err
		}
	}

	return &ntm, nil
}

// Render template
func (ntm notificationTemplate) Execute(templateType, templateName string, data interface{}) (string, error) {
	var (
		t  *template.Template
		ok bool
	)
	switch templateType {
	case entity.TemplateTypePushMessage:
		t, ok = ntm.pushMessageTemplate[templateName]
	case entity.TemplateTypeSMS:
		t, ok = ntm.smsTemplate[templateName]
	case entity.TemplateTypeEmail:
		t, ok = ntm.emailTemplate[templateName]
	}

	if !ok {
		return "", fmt.Errorf("templategen: template with type %s and name %s not found", templateType, templateName)
	}

	b := []byte{}
	buffer := bytes.NewBuffer(b)
	if err := t.Execute(buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// isExtSupported to check whether the file extension is valid for notification template
func (ntm notificationTemplate) isExtSupported(ext string) error {
	switch ext {
	case ".json":
	default:
		return errors.New("templategen: template format not valid, expecting json")
	}

	return nil
}

// TemplateSMS struct for loading sms notification template
type TemplateSMS struct {
	Message string `json:"message"`
}

func (ntm *notificationTemplate) loadSMS(templatePath string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// skip if  is dir
	if info.IsDir() {
		return nil
	}

	// now only support json
	if err := ntm.isExtSupported(info.Name()); err != nil {
		return err
	}

	out, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return err
	}

	smsTemplates := make(map[string]TemplateSMS)
	if err := json.Unmarshal(out, &smsTemplates); err != nil {
		return err
	}

	for templateName, smsTemplate := range smsTemplates {
		var err error

		t := template.New(templateName)
		t, err = t.Parse(smsTemplate.Message)
		if err != nil {
			return err
		}
	}

	return nil
}
