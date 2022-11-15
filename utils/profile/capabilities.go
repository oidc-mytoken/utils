package profile

import (
	"encoding/json"
	"strings"

	"github.com/oidc-mytoken/api/v0"
	"github.com/pkg/errors"

	"github.com/oidc-mytoken/utils/utils"
	"github.com/oidc-mytoken/utils/utils/jsonutils"
)

// ParseCapabilityTemplate parses the content of a capability template
func (p ProfileParser) ParseCapabilityTemplate(content []byte) (api.Capabilities, error) {
	capStrings, err := p.ParseCapabilityTemplateToStrings(content)
	capStrings = utils.UniqueSlice(capStrings)
	var caps api.Capabilities = nil
	if err == nil {
		caps = api.NewCapabilities(capStrings)
	}
	return caps, err
}

// ParseCapabilityTemplateToStringsByName parses the content of a capability template into a slice of strings
func (p ProfileParser) ParseCapabilityTemplateToStringsByName(name string) ([]string, error) {
	content, err := p.reader.ReadCapabilityTemplate(normalizeTemplateName(name))
	if err != nil {
		return nil, err
	}
	return p.ParseCapabilityTemplateToStrings(content)
}

// ParseCapabilityTemplateToStrings parses the content of a capability template into a slice of strings
func (p ProfileParser) ParseCapabilityTemplateToStrings(content []byte) (capStrings []string, err error) {
	if len(content) == 0 {
		return nil, nil
	}
	var tmpCapStrings []string
	if jsonutils.IsJSONArray(content) {
		if err = errors.WithStack(json.Unmarshal(content, &tmpCapStrings)); err != nil {
			return
		}
	} else {
		tmpCapStrings = strings.Split(string(content), " ")
	}
	for _, c := range tmpCapStrings {
		if !strings.HasPrefix(c, "@") {
			capStrings = append(capStrings, c)
		} else {
			templateCaps, e := p.ParseCapabilityTemplateToStringsByName(c[1:])
			if e != nil {
				err = e
				return
			}
			if len(templateCaps) > 0 {
				capStrings = append(capStrings, templateCaps...)
			}
		}
	}
	return
}
