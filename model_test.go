package mecab

import "testing"

func TestModel(t *testing.T) {
	model, err := NewModel(map[string]string{"output-format-type": "wakati"})
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
