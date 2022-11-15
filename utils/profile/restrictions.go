package profile

import (
	"encoding/json"
	"reflect"

	"github.com/oidc-mytoken/api/v0"
	"github.com/pkg/errors"

	"github.com/oidc-mytoken/utils/utils/jsonutils"
	"github.com/oidc-mytoken/utils/utils/timerestriction"
)

func parseRestrictionsTemplateByName(name string) (api.Restrictions, error) {
	content, err := templateReader.readRestrictionsTemplate(normalizeTemplateName(name))
	if err != nil {
		return nil, err
	}
	return ParseRestrictionsTemplate(content)
}

// ParseRestrictionsTemplate parses the content of a restrictions template
func ParseRestrictionsTemplate(content []byte) (api.Restrictions, error) {
	if len(content) == 0 {
		return nil, nil
	}

	var err error
	var restr []timerestriction.APIRestriction
	content, err = createFinalTemplate(content, templateReader.readRestrictionsTemplate)
	if err != nil {
		return nil, err
	}
	if len(content) == 0 {
		return nil, nil
	}
	if jsonutils.IsJSONObject(content) {
		content = jsonutils.Arrayify(content)
	}
	if err = errors.WithStack(json.Unmarshal(content, &restr)); err != nil {
		return nil, err
	}
	finalRestrs := make(api.Restrictions, 0)
	for _, r := range restr {
		if !reflect.DeepEqual(r, timerestriction.APIRestriction{}) {
			ar := api.Restriction(r)
			finalRestrs = append(finalRestrs, &ar)
		}
	}
	return finalRestrs, nil
}
