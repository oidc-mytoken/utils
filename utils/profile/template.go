package profile

import (
	"encoding/json"
	"path"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/pkg/errors"

	"github.com/oidc-mytoken/utils/utils/fileutil"
	"github.com/oidc-mytoken/utils/utils/jsonutils"
)

type readFnc func(string) ([]byte, error)

type TemplateReader interface {
	ReadProfile(string) ([]byte, error)
	ReadRestrictionsTemplate(string) ([]byte, error)
	ReadRotationTemplate(string) ([]byte, error)
	ReadCapabilityTemplate(string) ([]byte, error)
}

type FileTemplateReader struct {
	globalBaseDir string
	userBaseDir   string
	reader        readFnc
}

// ProfileParser is a type for parsing profiles and templates by using a TemplateReader
type ProfileParser struct {
	reader TemplateReader
}

// NewProfileParser creates a new ProfileParser from a TemplateReader
func NewProfileParser(reader TemplateReader) *ProfileParser {
	return &ProfileParser{reader: reader}
}

func init() {
	DefaultFileBasedProfileParser = NewProfileParser(newFileTemplateReader(fileutil.ReadFile))
}

// DefaultFileBasedProfileParser is the default profile parser for file-based parsing
var DefaultFileBasedProfileParser *ProfileParser

func newFileTemplateReader(reader readFnc) *FileTemplateReader {
	return &FileTemplateReader{
		globalBaseDir: "/etc/mytoken",
		userBaseDir:   userBasePath(),
		reader:        reader,
	}
}

func userBasePath() string {
	const conf = "~/.config/mytoken"
	const dot = "~/.mytoken"
	if fileutil.FileExists(conf) {
		return conf
	}
	return dot
}

// ReadFile reads any template file
func (r FileTemplateReader) ReadFile(relPath string) ([]byte, error) {
	globalP := path.Join(r.globalBaseDir, relPath)
	userP := path.Join(r.userBaseDir, relPath)
	global, _ := r.reader(globalP)
	user, _ := r.reader(userP)
	if len(user) == 0 {
		return global, nil
	}
	if len(global) == 0 {
		return user, nil
	}
	combined, _ := jsonpatch.MergePatch(global, user)
	return combined, nil
}

// ReadRestrictionsTemplate implements the TemplateReader interface
func (r FileTemplateReader) ReadRestrictionsTemplate(name string) ([]byte, error) {
	p := path.Join("restrictions.d", name)
	return r.ReadFile(p)
}

// ReadCapabilityTemplate implements the TemplateReader interface
func (r FileTemplateReader) ReadCapabilityTemplate(name string) ([]byte, error) {
	p := path.Join("capabilities.d", name)
	return r.ReadFile(p)
}

// ReadRotationTemplate implements the TemplateReader interface
func (r FileTemplateReader) ReadRotationTemplate(name string) ([]byte, error) {
	p := path.Join("rotation.d", name)
	return r.ReadFile(p)
}

// ReadProfile implements the TemplateReader interface
func (r FileTemplateReader) ReadProfile(name string) ([]byte, error) {
	p := path.Join("profiles.d", name)
	return r.ReadFile(p)
}

func normalizeTemplateName(name string) string {
	if name != "" && name[0] == '@' {
		return name[1:]
	}
	return name
}

type includeTemplates struct {
	Include json.RawMessage `json:"include"`
}

func createFinalTemplate(content []byte, read readFnc) ([]byte, error) {
	if len(content) == 0 {
		return nil, nil
	}
	baseIsArray := jsonutils.IsJSONArray(content)
	if baseIsArray {
		var contents []json.RawMessage
		if err := errors.WithStack(json.Unmarshal(content, &contents)); err != nil {
			return nil, err
		}
		final := []byte(`[]`)
		for _, c := range contents {
			c = jsonutils.UnwrapString(c)
			cf, err := createFinalTemplate(c, read)
			if err != nil {
				return nil, err
			}
			if !jsonutils.IsJSONArray(cf) {
				cf = jsonutils.Arrayify(cf)
			}
			final = jsonutils.MergeJSONArrays(final, cf)
		}
		return final, nil
	}

	if !jsonutils.IsJSONObject(content) {
		// must be one or multiple template names
		templates := strings.Split(string(content), " ")
		if len(templates) == 1 {
			// must be single template name
			c, err := read(normalizeTemplateName(templates[0]))
			if err != nil {
				return nil, err
			}
			return createFinalTemplate(c, read)
		}
		// multiple templates
		return parseIncludes([]byte(`{}`), templates, read)
	}
	var inc includeTemplates
	if err := errors.WithStack(json.Unmarshal(content, &inc)); err != nil {
		return nil, err
	}
	includes := make([]string, 0)
	if len(inc.Include) > 0 {
		if inc.Include[0] == '[' {
			if err := json.Unmarshal(inc.Include, &includes); err != nil {
				return nil, err
			}
		} else {
			inc.Include = jsonutils.UnwrapString(inc.Include)
			includes = strings.Split(string(inc.Include), " ")
		}
	}
	return parseIncludes(content, includes, read)
}

func parseIncludes(content []byte, includes []string, read readFnc) ([]byte, error) {
	baseIsArray := jsonutils.IsJSONArray(content)
	for _, inP := range includes {
		c, err := read(normalizeTemplateName(inP))
		if err != nil {
			return nil, err
		}
		includeContent, err := createFinalTemplate(c, read)
		if err != nil {
			return nil, err
		}
		isArray := jsonutils.IsJSONArray(includeContent)
		if !baseIsArray && !isArray {
			content, _ = jsonpatch.MergePatch(includeContent, content)
			continue
		}
		if !baseIsArray {
			content = jsonutils.Arrayify(content)
			baseIsArray = true
		}
		if !isArray {
			includeContent = jsonutils.Arrayify(includeContent)
		}
		content = jsonutils.MergeJSONArrays(content, includeContent)
	}
	return content, nil
}
