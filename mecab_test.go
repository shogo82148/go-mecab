package mecab

import "testing"

func TestMeCab(t *testing.T) {
	m, err := New(map[string]string{"output-format-type": "wakati"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	defer m.Destroy()

	tagger, err := m.NewTagger()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	defer tagger.Destroy()

	result, err := tagger.Parse("こんにちは世界")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	expected := "こんにちは 世界 \n"
	if result != expected {
		t.Errorf("want `%s`, but `%s`", expected, result)
	}
}
