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
    echo "CGO_LDFLAGS=-L$(cygpath -w /mingw64/lib) -lmecab -lstdc++"
    echo "CGO_CFLAGS=-I$(cygpath -w /mingw64/include)"
} >> "$GITHUB_ENV"

cat << "DIC" > /mingw64/lib/mecab/dic/ipadic/dicrc
;
; Configuration file of IPADIC
;
; $Id: dicrc,v 1.4 2006/04/08 06:41:36 taku-ku Exp $;
;

dicdir = $(cygpath -w /mingw64/lib/mecab/dic/ipadic)

cost-factor = 800
bos-feature = BOS/EOS,*,*,*,*,*,*,*,*
eval-size = 8
unk-eval-size = 4
config-charset = EUC-JP

; yomi
node-format-yomi = %pS%f[7]
unk-format-yomi = %M
eos-format-yomi  = \\n
DIC

# The default mecabrc path is "C:\Program Files\mecab\etc\mecabrc" if mecab is built with mingw32-w64.
# but it is not correct in MSYS2 environment.
echo "MECABRC_PATH=$(cygpath -w /mingw64/etc/mecabrc)" >> "$GITHUB_ENV"
cygpath -w /mingw64/bin >> "$GITHUB_PATH"
