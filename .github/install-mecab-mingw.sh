#!/bin/bash
# from https://gist.github.com/dtan4/351d031bec0c3d45cd8f
# see also http://qiita.com/dtan4/items/c6a087666296fbd5fffb

set -uxe

TMPDIR=$(mktemp -d)
trap 'rm -rfv "$TMPDIR"' EXIT

MECAB_VERSION=0.996.2
IPADIC_VERSION=2.7.0-20070801
# install mecab
cd "$TMPDIR"
curl -o mecab.tar.gz -sSL "https://github.com/shogo82148/mecab/releases/download/v$MECAB_VERSION/mecab-$MECAB_VERSION.tar.gz"
tar zxfv mecab.tar.gz
cd "mecab-$MECAB_VERSION"
./configure --enable-utf8-only
make -j2 LDFLAGS=-static
make check
sudo make install
sudo ldconfig

cd "$TMPDIR"
curl -o mecab-ipadic.tar.gz -sSL "https://github.com/shogo82148/mecab/releases/download/v$MECAB_VERSION/mecab-ipadic-$IPADIC_VERSION.tar.gz"
tar zxfv mecab-ipadic.tar.gz
cd "mecab-ipadic-$IPADIC_VERSION"
./configure --with-charset=utf8
make
sudo make install

echo "CGO_LDFLAGS=$(mecab-config --libs)" >> "$GITHUB_ENV"
echo "CGO_CFLAGS=-I$(mecab-config --inc-dir)" >> "$GITHUB_ENV"
