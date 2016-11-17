FROM golang:latest
MAINTAINER Mephis Pheies <mephistommm@gmail.com>

ENV REBUILD_DAY 02.10.2016

# modify localtime
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# change source.list
COPY ./sources.list /etc/apt/sources.list

# set gobin env
ENV GOBIN=/go/bin

# install unzip 
RUN apt-get update && apt-get install -y unzip
