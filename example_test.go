package mecab_test

import (
	"fmt"

	"github.com/shogo82148/go-mecab"
)

func ExampleMeCab_Parse() {
	tagger, err := mecab.New(map[string]string{})
	if err != nil {
		panic(err)
	}
	defer tagger.Destroy()

	result, err := tagger.Parse("こんにちは世界")
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
	// Output:
	// こんにちは	感動詞,*,*,*,*,*,こんにちは,コンニチハ,コンニチワ
	// 世界	名詞,一般,*,*,*,*,世界,セカイ,セカイ
	// EOS
}

func ExampleMeCab_ParseLattice() {
	tagger, err := mecab.New(map[string]string{})
	if err != nil {
		panic(err)
	}
	defer tagger.Destroy()

	lattice, err := mecab.NewLattice()
	if err != nil {
		panic(err)
	}

	lattice.SetSentence("こんにちは世界")
	err = tagger.ParseLattice(lattice)
	if err != nil {
		panic(err)
	}
	fmt.Println(lattice.String())
	// Output:
	// こんにちは	感動詞,*,*,*,*,*,こんにちは,コンニチハ,コンニチワ
	// 世界	名詞,一般,*,*,*,*,世界,セカイ,セカイ
	// EOS
}

func ExampleMeCab_ParseToNode() {
	tagger, err := mecab.New(map[string]string{})
	if err != nil {
		panic(err)
	}
	defer tagger.Destroy()

	// XXX: avoid GC problem with MeCab 0.996 (see https://github.com/taku910/mecab/pull/24)
	tagger.Parse("")

	node, err := tagger.ParseToNode("こんにちは世界")
	if err != nil {
		panic(err)
	}

	for ; !node.IsZero(); node = node.Next() {
		fmt.Printf("%s\t%s\n", node.Surface(), node.Feature())
	}
	// Output:
	// 	BOS/EOS,*,*,*,*,*,*,*,*
	// こんにちは	感動詞,*,*,*,*,*,こんにちは,コンニチハ,コンニチワ
	// 世界	名詞,一般,*,*,*,*,世界,セカイ,セカイ
	// 	BOS/EOS,*,*,*,*,*,*,*,*
}
