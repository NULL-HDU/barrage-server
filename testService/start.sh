#!/bin/env bash 
#This file start the test service after installing dependences and 
#install server self.
#
#Author: Mephis Pheies <mephistommm@gmail.com>

go get -u golang.org/x/net/websocket
go install barrage-test-server.go

barrage-test-server
