package timerestriction

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/oidc-mytoken/api/v0"
	"github.com/pkg/errors"

	"github.com/oidc-mytoken/utils/unixtime"
	"github.com/oidc-mytoken/utils/utils/duration"
)

// ParseTime parses a time string
func ParseTime(t string) (int64, error) {
	if t == "" {
		return 0, nil
	}
	i, err := strconv.ParseInt(t, 10, 64)
	if err == nil {
		if t[0] == '+' {
			return int64(unixtime.InSeconds(i)), nil
		}
		return i, nil
	}
	if t[0] == '+' {
		d, err := duration.ParseDuration(t[1:])
		return int64(unixtime.New(time.Now().Add(d))), err
	}
	tt, err := time.ParseInLocation("2006-01-02 15:04", t, time.Local)
	return int64(unixtime.New(tt)), err
}

type restrictionWT struct {
	api.Restriction
	ExpiresAt timeValue `json:"exp"`
	NotBefore timeValue `json:"nbf"`
}

type timeValue struct {
	Value int64
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (tv *timeValue) UnmarshalJSON(data []byte) error {
	if err := errors.WithStack(json.Unmarshal(data, &tv.Value)); err == nil {
		return nil
	}
	var str string
	if err := errors.WithStack(json.Unmarshal(data, &str)); err != nil {
		return err
	}
	t, err := ParseTime(str)
	if err != nil {
		return err
	}
	tv.Value = t
	return nil
}

// APIRestriction is a type for extended an api.Restriction
type APIRestriction api.Restriction

// UnmarshalJSON implements the json.Unmarshaler interface
func (r *APIRestriction) UnmarshalJSON(data []byte) error {
	rr := restrictionWT{}
	if err := json.Unmarshal(data, &rr); err != nil {
		return err
	}
	rr.Restriction.ExpiresAt = rr.ExpiresAt.Value
	rr.Restriction.NotBefore = rr.NotBefore.Value
	*r = APIRestriction(rr.Restriction)
	return nil
}
