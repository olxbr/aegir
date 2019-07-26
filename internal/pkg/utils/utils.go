package utils

import "strings"

type jsonPatchOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

type Violation struct {
	RuleName     string                 `json:"rule_name"`
	Description  string                 `json:"description"`
	JSONPath     string                 `json:"json_path"`
	Object       map[string]interface{} `json:"object"`
	Message      string                 `json:"error,omitempty"`
	SlackChannel string                 `json:"slack_channel,omitempty"`
}

func GetLastField(field string) string {
	s := strings.Split(strings.TrimRight(field, "/"), ".")
	return s[len(s)-1]
}
