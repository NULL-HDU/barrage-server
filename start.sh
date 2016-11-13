#!/bin/env bash 
#This file start the barrage server after installing dependences and 
#install server self.
#
#Author: Mephis Pheies <mephistommm@gmail.com>

go get -u golang.org/x/net/websocket
go install main.go

barrage-server
