package profile

import (
	"encoding/json"

	"github.com/oidc-mytoken/api/v0"
	"github.com/pkg/errors"
)

type profileUnmarshal struct {
	api.GeneralMytokenRequest
	Restrictions         jsonraw `json:"restrictions"`
	Capabilities         jsonraw `json:"capabilities"`
	SubtokenCapabilities jsonraw `json:"subtoken_capabilities"`
	Rotation             jsonraw `json:"rotation"`
}

type jsonraw string

func (r *jsonraw) UnmarshalJSON(data []byte) error {
	var raw json.RawMessage
	if err := errors.WithStack(json.Unmarshal(data, &raw)); err != nil {
		return err
	}
	rawStr := string(raw)
	if rawStr != "" && rawStr[0] == '"' && rawStr[len(rawStr)-1] == '"' {
		rawStr = rawStr[1 : len(rawStr)-1]
	}
	*r = jsonraw(rawStr)
	return nil
}

// ParseProfileByName parses the content of a profile by name
func (p ProfileParser) ParseProfileByName(name string) (api.GeneralMytokenRequest, error) {
	content, err := p.reader.ReadProfile(normalizeTemplateName(name))
	if err != nil {
		return api.GeneralMytokenRequest{}, err
	}
	return p.ParseProfile(content)
}

// ParseProfile parses the content of a profile
func (p ProfileParser) ParseProfile(content []byte) (api.GeneralMytokenRequest, error) {
	if len(content) == 0 {
		return api.GeneralMytokenRequest{}, nil
	}
	var err error
	var pj profileUnmarshal
	content, err = createFinalTemplate(content, p.reader.ReadProfile)
	if err != nil {
		return pj.GeneralMytokenRequest, err
	}
	if len(content) > 0 {
		if err = errors.WithStack(json.Unmarshal(content, &pj)); err != nil {
			return pj.GeneralMytokenRequest, err
		}
	}
	pj.GeneralMytokenRequest.Rotation, err = p.ParseRotationTemplate([]byte(pj.Rotation))
	if err != nil {
		return pj.GeneralMytokenRequest, err
	}
	pj.GeneralMytokenRequest.Capabilities, err = p.ParseCapabilityTemplate([]byte(pj.Capabilities))
	if err != nil {
		return pj.GeneralMytokenRequest, err
	}
	pj.GeneralMytokenRequest.Restrictions, err = p.ParseRestrictionsTemplate([]byte(pj.Restrictions))
	if err != nil {
		return pj.GeneralMytokenRequest, err
	}
	return pj.GeneralMytokenRequest, nil
}
