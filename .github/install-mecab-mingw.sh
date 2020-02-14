#!/bin/bash

set -uxe

# setup MSYS2
pacman-key --init
pacman-key --populate msys2
pacman -Sy
pacman -S --verbose --noconfirm --noprogressbar --needed make mingw-w64-x86_64-gcc mingw-w64-x86_64-libtool

TMPDIR=$(mktemp -d)
trap 'rm -rfv "$TMPDIR"' EXIT

export PATH=/mingw64/bin:$PATH

MECAB_VERSION=0.996.2
IPADIC_VERSION=2.7.0-20070801
# install mecab
cd "$TMPDIR"
curl -o mecab.tar.gz -sSL "https://github.com/shogo82148/mecab/releases/download/v$MECAB_VERSION/mecab-$MECAB_VERSION.tar.gz"
tar zxfv mecab.tar.gz
cd "mecab-$MECAB_VERSION"
./configure --enable-utf8-only
make
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

echo "::set-env name=CGO_LDFLAGS::$(mecab-config --libs)"
echo "::set-env name=CGO_CFLAGS::-I$(mecab-config --inc-dir)"
