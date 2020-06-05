#!/bin/sh
rm liblz77.so
rm ../liblz77.so
g++ Lz77.cpp -fPIC -shared -o liblz77.so
cp liblz77.so ../
