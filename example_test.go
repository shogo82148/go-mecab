package mecab_test

import (
	"fmt"
	"os"

	"github.com/shogo82148/go-mecab"
)

func ExampleMeCab_Parse() {
	options := map[string]string{}
	if path := os.Getenv("MECABRC_PATH"); path != "" {
		options["rcfile"] = path
	}

	tagger, err := mecab.New(options)
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
	options := map[string]string{}
	if path := os.Getenv("MECABRC_PATH"); path != "" {
		options["rcfile"] = path
	}

	tagger, err := mecab.New(options)
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

func ExampleMeCab_ParseLattice_nBest() {
	options := map[string]string{}
	if path := os.Getenv("MECABRC_PATH"); path != "" {
		options["rcfile"] = path
	}

	tagger, err := mecab.New(options)
	if err != nil {
		panic(err)
	}
	defer tagger.Destroy()

	lattice, err := mecab.NewLattice()
	if err != nil {
		panic(err)
	}

	lattice.SetSentence("こんにちは世界")
	lattice.AddRequestType(mecab.RequestTypeNBest)
	err = tagger.ParseLattice(lattice)
	if err != nil {
		panic(err)
	}
	for i := 0; i < 5; i++ {
		fmt.Println(lattice.String())
		if !lattice.Next() {
			break
		}
	}
	// Output:
	// こんにちは	感動詞,*,*,*,*,*,こんにちは,コンニチハ,コンニチワ
	// 世界	名詞,一般,*,*,*,*,世界,セカイ,セカイ
	// EOS
	//
	// こんにちは	感動詞,*,*,*,*,*,こんにちは,コンニチハ,コンニチワ
	// 世界	名詞,一般,*,*,*,*,世界,セカイ,セカイ
	// EOS
	//
	// こんにちは	感動詞,*,*,*,*,*,こんにちは,コンニチハ,コンニチワ
	// 世	名詞,一般,*,*,*,*,世,ヨ,ヨ
	// 界	名詞,接尾,一般,*,*,*,界,カイ,カイ
	// EOS
	//
	// こんにちは	感動詞,*,*,*,*,*,こんにちは,コンニチハ,コンニチワ
	// 世	名詞,一般,*,*,*,*,世,ヨ,ヨ
	// 界	名詞,固有名詞,地域,一般,*,*,界,サカイ,サカイ
	// EOS
	//
	// こんにちは	感動詞,*,*,*,*,*,こんにちは,コンニチハ,コンニチワ
	// 世	名詞,接尾,助数詞,*,*,*,世,セイ,セイ
	// 界	名詞,接尾,一般,*,*,*,界,カイ,カイ
	// EOS
}

func ExampleMeCab_ParseToNode() {
	options := map[string]string{}
	if path := os.Getenv("MECABRC_PATH"); path != "" {
		options["rcfile"] = path
	}

	tagger, err := mecab.New(options)
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
