#!/bin/bash
# from https://gist.github.com/dtan4/351d031bec0c3d45cd8f
# see also http://qiita.com/dtan4/items/c6a087666296fbd5fffb

set -uxe

TMPDIR=$(mktemp -d)
trap 'rm -rfv "$TMPDIR"' EXIT

MECAB_VERSION=0.996.4
IPADIC_VERSION=2.7.0-20070801
# install mecab
cd "$TMPDIR"
curl -o mecab.tar.gz -sSL "https://github.com/shogo82148/mecab/releases/download/v$MECAB_VERSION/mecab-$MECAB_VERSION.tar.gz"
tar zxfv mecab.tar.gz
cd "mecab-$MECAB_VERSION"
./configure --enable-utf8-only --host=x86_64-w64-mingw32
make -j2
# make check # it fails :(
make install

cd "$TMPDIR"
curl -o mecab-ipadic.tar.gz -sSL "https://github.com/shogo82148/mecab/releases/download/v$MECAB_VERSION/mecab-ipadic-$IPADIC_VERSION.tar.gz"
tar zxfv mecab-ipadic.tar.gz
cd "mecab-ipadic-$IPADIC_VERSION"
./configure --with-charset=utf8
make
make install

{
    echo "CGO_LDFLAGS=C:\\msys64\\lib"
    echo "CGO_CFLAGS=-IC:\\msys64\\include"
} >> "$GITHUB_ENV"

# The default mecabrc path is "C:\Program Files\mecab\etc\mecabrc" if mecab is built with mingw32-w64.
# but it is not correct in MSYS2 environment.
echo "MECABRC_PATH=C:\msys64\mingw64\etc\mecabrc" >> "$GITHUB_ENV"
echo "C:\msys64\mingw64\bin" >> "$GITHUB_PATH"
