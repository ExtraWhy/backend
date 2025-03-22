#!/bin/bash

if [[ -z $1 ]]; then 
	echo 'Perovide golang download link' && exit
fi

wget $1

tar -C /usr/local -xzf *.gz

rm -rf /usr/local/go

cp -rvf go/ /usr/local


echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc 

go version