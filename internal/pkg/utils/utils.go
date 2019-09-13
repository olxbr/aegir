package utils

import (
	"os"
	"strings"
)

type Violation struct {
	RuleName     string                 `json:"rule_name"`
	Description  string                 `json:"description"`
	JSONPath     string                 `json:"json_path"`
	Object       map[string]interface{} `json:"object"`
	Message      string                 `json:"error,omitempty"`
	SlackChannel string                 `json:"slack_channel,omitempty"`
}

//GetLastField returns the last word of a path delimited by '/'
func GetLastField(field string) string {
	s := strings.Split(strings.TrimRight(field, "/"), ".")
	return s[len(s)-1]
}

func index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

//Include returns true if string is in slice and false otherwise
func Include(vs []string, t string) bool {
	return index(vs, t) >= 0
}

func GetEnvAsSlice(name string, sep string) ([]string, bool) {
	valStr, ok := os.LookupEnv(name)
	return strings.Split(valStr, sep), ok
}
