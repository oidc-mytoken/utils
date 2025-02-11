package jwtutils

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/lestrrat-go/jwx/jwt"
	log "github.com/sirupsen/logrus"
)

// GetFromJWT returns the values for the requested keys from the JWT (or all)
func GetFromJWT(rlog log.Ext1FieldLogger, token string, key ...string) map[string]interface{} {
	tokenJWT, err := jwt.Parse([]byte(token))
	if err != nil || tokenJWT == nil {
		return nil
	}
	rlog.Trace("Parsed token")
	claims, err := tokenJWT.AsMap(context.Background())
	if err != nil {
		return nil
	}
	if len(key) == 0 {
		return claims
	}
	values := make(map[string]interface{}, len(key))
	for _, k := range key {
		v, set := claims[k]
		if set {
			values[k] = v
		}
	}
	return values
}

// GetValueFromJWT returns the value for the given key
func GetValueFromJWT(rlog log.Ext1FieldLogger, token, key string) interface{} {
	res := GetFromJWT(rlog, token)
	return res[key]
}

// GetStringFromJWT returns a string value for the given key
func GetStringFromJWT(rlog log.Ext1FieldLogger, token, key string) (string, bool) {
	res := GetFromJWT(rlog, token)
	v, ok := res[key]
	if !ok {
		return "", false
	}
	vv, ok := v.(string)
	return vv, ok
}

// GetAudiencesFromJWT parses the passed jwt token and returns the aud claim as a slice of strings
func GetAudiencesFromJWT(rlog log.Ext1FieldLogger, token string) ([]string, bool) {
	rlog.Trace("Getting auds from token")
	res := GetFromJWT(rlog, token)
	auds, found := res["aud"]
	if !found {
		return nil, false
	}
	switch v := auds.(type) {
	case string:
		return []string{v}, true
	case []string:
		return v, true
	case []interface{}:
		strs := []string{}
		for _, s := range v {
			str, ok := s.(string)
			if !ok {
				return nil, false
			}
			strs = append(strs, str)
		}
		return strs, true
	default:
		return nil, false
	}
}

// IsJWT checks if a string is a jwt
func IsJWT(token string) bool {
	arr := strings.Split(token, ".")
	if len(arr) < 3 {
		return false
	}
	for i, segment := range arr {
		if segment != "" || i < 2 { // first two segments must not be empty
			if _, err := base64.URLEncoding.DecodeString(arr[2]); err != nil {
				return false
			}
		}
	}
	return true
}
