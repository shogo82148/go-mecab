# go-mecab

[![Build Status](https://travis-ci.org/shogo82148/go-mecab.svg?branch=master)](https://travis-ci.org/shogo82148/go-mecab)

go-mecab is [MeCab](http://taku910.github.io/mecab/) binding for Golang.

## SYNOPSIS

``` go
import "github.com/shogo82148/go-mecab"

tagger, err := mecab.New(map[string]string{"output-format-type": "wakati"})
defer tagger.Destroy()
result, err := tagger.Parse("こんにちは世界")
fmt.Println(result)
// Output: こんにちは 世界
```

## INSTALL

You need to tell Go where MeCab has been installed.

``` bash
$ export CGO_LDFLAGS="-L/path/to/lib -lmecab -lstdc++"
$ export CGO_CFLAGS="-I/path/to/include"
$ go get github.com/shogo82148/go-mecab
```

If you installed `mecab-config`, execute following comands.

``` bash
$ export CGO_LDFLAGS="`mecab-config --libs`"
$ export CGO_FLAGS="`mecab-config --inc-dir`"
$ go get github.com/shogo82148/go-mecab
```

## SEE ALSO

- [godoc](https://godoc.org/github.com/shogo82148/go-mecab)
- [MeCab](http://taku910.github.io/mecab/)
- [MeCab repository](https://github.com/taku910/mecab)
