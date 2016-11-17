#!/bin/env bash 
#This file start the barrage server after installing dependences and 
#install server self.
#
#Author: Mephis Pheies <mephistommm@gmail.com>

NET_PKG="net-4971afdc2f162e82d185353533d3cf16188a9f4e"

if [ ! -d /go/src/golang.org/ ]; then
    wget http://gopm.dn.qbox.me/golang.org/x/$NET_PKG.zip -O /go/src/net.zip
    mkdir -p /go/src/golang.org/x
    unzip /go/src/net.zip -d /go/src/golang.org/x/
    mv /go/src/golang.org/x/$NET_PKG /go/src/golang.org/x/net
    rm /go/src/net.zip
fi

go install

barrage-server
