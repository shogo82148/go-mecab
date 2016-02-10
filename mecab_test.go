package mecab

import "testing"

func TestMeCab(t *testing.T) {
	m, _ := New(map[string]string{})
	defer m.Destroy()
}
