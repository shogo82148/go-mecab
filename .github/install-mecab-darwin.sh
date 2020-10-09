#!/bin/bash
# from https://gist.github.com/dtan4/351d031bec0c3d45cd8f
# see also http://qiita.com/dtan4/items/c6a087666296fbd5fffb

set -uxe

TMPDIR=$(mktemp -d)
trap 'rm -rfv "$TMPDIR"' EXIT

PREFIX=$GITHUB_WORKSPACE/mecab
export PATH=$PREFIX/bin:$PATH

if [ ! -e "$PREFIX/bin/mecab" ]; then
    # install mecab
    cd "$TMPDIR"
    curl -o mecab.tar.gz -sSL "https://github.com/shogo82148/mecab/releases/download/v$MECAB_VERSION/mecab-$MECAB_VERSION.tar.gz"
    tar zxfv mecab.tar.gz
    cd "mecab-$MECAB_VERSION"
    ./configure --enable-utf8-only --prefix="$PREFIX"
    make -j2
    make check
    sudo make install

    cd "$TMPDIR"
    curl -o mecab-ipadic.tar.gz -sSL "https://github.com/shogo82148/mecab/releases/download/v$MECAB_VERSION/mecab-ipadic-$IPADIC_VERSION.tar.gz"
    tar zxfv mecab-ipadic.tar.gz
    cd "mecab-ipadic-$IPADIC_VERSION"
    ./configure --with-charset=utf8 --prefix="$PREFIX"
    make
    sudo make install
fi

cat << CGO_FLAGS >> "$GITHUB_ENV"
LD_LIBRARY_PATH=$PREFIX/lib
CGO_LDFLAGS=$(mecab-config --libs)
CGO_CFLAGS=-I$(mecab-config --inc-dir)
CGO_FLAGS
