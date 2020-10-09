#!/bin/bash
# from https://gist.github.com/dtan4/351d031bec0c3d45cd8f
# see also http://qiita.com/dtan4/items/c6a087666296fbd5fffb

set -uxe

TMPDIR=$(mktemp -d)
trap 'rm -rfv "$TMPDIR"' EXIT

PREFIX=$(cygpath -u "$GITHUB_WORKSPACE\mecab")
export PATH=$PREFIX/bin:$PATH

if [ ! -e "$PREFIX/bin/mecab" ]; then
    # install mecab
    cd "$TMPDIR"
    curl -o mecab.tar.gz -sSL "https://github.com/shogo82148/mecab/releases/download/v$MECAB_VERSION/mecab-$MECAB_VERSION.tar.gz"
    tar zxfv mecab.tar.gz
    cd "mecab-$MECAB_VERSION"
    ./configure --enable-utf8-only --host=x86_64-w64-mingw32 --prefix="$PREFIX"
    make -j2
    # make check # it fails :(
    make install

    cd "$TMPDIR"
    curl -o mecab-ipadic.tar.gz -sSL "https://github.com/shogo82148/mecab/releases/download/v$MECAB_VERSION/mecab-ipadic-$IPADIC_VERSION.tar.gz"
    tar zxfv mecab-ipadic.tar.gz
    cd "mecab-ipadic-$IPADIC_VERSION"
    ./configure --with-charset=utf8 --prefix="$PREFIX"
    make
    make install
fi

{
    echo "CGO_LDFLAGS=-L$(cygpath -w /mingw64/lib) -L$(cygpath -w "$PREFIX/lib") -lmecab -lstdc++"
    echo "CGO_CFLAGS=-I$(cygpath -w /mingw64/include) -I$(cygpath -w "$PREFIX/include")"
} >> "$GITHUB_ENV"

cat << DIC > "$PREFIX/etc/mecabrc"
dicdir = $(cygpath -w "$PREFIX/lib/mecab/dic/ipadic")
DIC

# The default mecabrc path is "C:\Program Files\mecab\etc\mecabrc" if mecab is built with mingw32-w64.
# but it is not correct in MSYS2 environment.
echo "MECABRC_PATH=$(cygpath -w "$PREFIX/etc/mecabrc")" >> "$GITHUB_ENV"
cygpath -w /mingw64/bin >> "$GITHUB_PATH"
cygpath -w "$PREFIX/bin" >> "$GITHUB_PATH"
