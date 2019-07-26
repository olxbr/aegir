package utils

import "testing"

func TestGetLastField(t *testing.T) {
	lastfield := GetLastField("metadata.spec.template")
	expected := "template"

	if lastfield != expected {
		t.Errorf("expected '%s' but got '%s'", expected, lastfield)
	}
}
