package profile

import (
	"encoding/json"

	"github.com/oidc-mytoken/api/v0"
	"github.com/pkg/errors"
)

// ParseRotationTemplateByName parses the content of a rotation template by name
func (p ProfileParser) ParseRotationTemplateByName(name string) (*api.Rotation, error) {
	content, err := p.reader.ReadRotationTemplate(normalizeTemplateName(name))
	if err != nil {
		return nil, err
	}
	return p.ParseRotationTemplate(content)
}

// ParseRotationTemplate parses the content of a rotation template
func (p ProfileParser) ParseRotationTemplate(content []byte) (*api.Rotation, error) {
	if len(content) == 0 {
		return nil, nil
	}
	var err error
	var rot api.Rotation
	content, err = createFinalTemplate(content, p.reader.ReadRotationTemplate)
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, nil
	}
	err = errors.WithStack(json.Unmarshal(content, &rot))
	return &rot, err
}
