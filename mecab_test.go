package mecab

import "testing"

func TestNewMeCab(t *testing.T) {
	mecab, err := New(map[string]string{"output-format-type": "wakati", "all-morphs": ""})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	defer mecab.Destroy()
}

func TestParse(t *testing.T) {
	mecab, err := New(map[string]string{"output-format-type": "wakati"})
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

func TestParseToNode(t *testing.T) {
	mecab, err := New(map[string]string{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	defer mecab.Destroy()

	// XXX: avoid GC, MeCab 0.996 has GC problem (see https://github.com/taku910/mecab/pull/24)
	mecab.Parse("")

	node, err := mecab.ParseToNode("こんにちは世界")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	node = node.Next()
	if node.Surface() != "こんにちは" {
		t.Errorf("want こんにちは, but %s", node.Surface())
	}
	node = node.Next()
	if node.Surface() != "世界" {
		t.Errorf("want 世界, but %s", node.Surface())
	}
}
