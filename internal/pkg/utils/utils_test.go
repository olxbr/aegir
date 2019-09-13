package utils

import (
	"os"
	"reflect"
	"testing"
)

func TestGetLastField(t *testing.T) {
	lastfield := GetLastField("metadata.spec.template")
	expected := "template"

	if lastfield != expected {
		t.Errorf("expected '%s' but got '%s'", expected, lastfield)
	}
}

func TestIndex(t *testing.T) {
	expected := 1
	strgs := []string{"hello", "world", "computer"}
	result := index(strgs, "world")
	if result != expected {
		t.Errorf("expected '%d' but got '%d'", expected, result)
	}
}

func TestIndexNotFound(t *testing.T) {
	expected := -1
	strgs := []string{"hello", "world", "computer"}
	result := index(strgs, "mundo")
	if result != expected {
		t.Errorf("expected '%d' but got '%d'", expected, result)
	}
}

func TestInclude(t *testing.T) {
	expected := true
	strgs := []string{"hello", "world", "computer"}
	result := Include(strgs, "hello")
	if result != expected {
		t.Errorf("expected '%t' but got '%t'", expected, result)
	}
}

func TestGetEnvAsSlice(t *testing.T) {
	expected := []string{"hello", "world", "computer"}
	os.Setenv("SLICE", "hello,world,computer")
	result, _ := GetEnvAsSlice("SLICE", ",")
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expected '%s' but got '%s'", expected, result)
	}
}
