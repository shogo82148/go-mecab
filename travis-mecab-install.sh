#!/bin/bash
# from https://gist.github.com/dtan4/351d031bec0c3d45cd8f
# see also http://qiita.com/dtan4/items/c6a087666296fbd5fffb

base_dir=`pwd`

wget https://mecab.googlecode.com/files/mecab-0.996.tar.gz
tar zxfv mecab-0.996.tar.gz
cd mecab-0.996
./configure --enable-utf8-only
make
make check
sudo make install
sudo ldconfig

cd $base_dir

wget https://mecab.googlecode.com/files/mecab-ipadic-2.7.0-20070801.tar.gz
tar zxfv mecab-ipadic-2.7.0-20070801.tar.gz
cd mecab-ipadic-2.7.0-20070801
./configure --with-charset=utf8
make
sudo make install
sudo ldconfig

cd $base_dir
rm -rf mecab-0.996 mecab-ipadic-2.7.0-20070801
