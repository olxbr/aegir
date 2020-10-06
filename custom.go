package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/grupozap/aegir/internal/pkg/utils"
	livr "github.com/k33nice/go-livr"
)

var (
	defaultRegex = regexp.MustCompile(".*")
)

func neq(args ...interface{}) livr.Validation {
	not_allowed := utils.FirstArg(args...)

	return func(value interface{}, builders ...interface{}) (interface{}, interface{}) {
		if value == nil || value == "" {
			return nil, nil
		}

		switch value.(type) {
		case float64, string, bool:
		default:
			return nil, errors.New("FORMAT_ERROR")
		}
		if fmt.Sprint(value) == fmt.Sprint(not_allowed) {
			return nil, errors.New("NOT_ALLOWED_VALUE")
		}

		return value, nil
	}
}

func not_like(args ...interface{}) livr.Validation {
	var re *regexp.Regexp
	var flags string
	if len(args) > 0 {
		if len(args) > 1 {
			if v, ok := args[1].(string); ok {
				if v == "i" {
					flags = "(?i)"
				}
			}
		}

		if v, ok := args[0].(string); ok {
			reg, err := regexp.Compile(flags + v)
			if err != nil {
				re = defaultRegex
			} else {
				re = reg
			}
		}
	}

	return func(value interface{}, builders ...interface{}) (interface{}, interface{}) {
		if value == nil || value == "" {
			return value, nil
		}

		switch v := value.(type) {
		case string:
			if matches := re.MatchString(v); matches {
				return nil, errors.New("WRONG_FORMAT")
			}
			return v, nil
		case float64:
			if matches := re.MatchString(strconv.FormatFloat(v, 'f', -1, 64)); matches {
				return nil, errors.New("WRONG_FORMAT")
			}
			return v, nil
		default:
			return nil, errors.New("FORMAT_ERROR")
		}
	}
}

var customRules map[string]livr.Builder

func init() {
	customRules = map[string]livr.Builder{
		"neq":      neq,
		"not_like": not_like,
	}

	v := livr.New(&livr.Options{})
	v.RegisterDefaultRules(customRules)
}
