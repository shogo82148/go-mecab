package mecab

import (
	"strings"
	"testing"
)

func TestModel(t *testing.T) {
	model, err := NewModel(rcfile(map[string]string{
		"output-format-type": "wakati",
	}))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	defer model.Destroy()

	mecab, err := model.NewMeCab()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	defer mecab.Destroy()

	result, err := mecab.Parse("こんにちは世界")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	expected := "こんにちは 世界 \n"
	if result != expected {
		t.Errorf("want `%s`, but `%s`", expected, result)
	}
}

func TestNewModel_error(t *testing.T) {
	_, err := NewModel(rcfile(map[string]string{
		"output-format-type": "unknown format",
	}))
	if err == nil {
		t.Errorf("expected error, but not")
		return
	}
	if !strings.Contains(err.Error(), "unknown format type [unknown format]") {
		t.Errorf("want %q error, got %q", "unknown format type [unknown format]", err.Error())
	}
}
