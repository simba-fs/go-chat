#!/bin/bash
mkdir -p bin
rm bin/* 
cd cmd
for i in $(ls);do
	echo building \"$i\"
	cd $i
	go build -o ../../bin
	cd ..
done
