#Dockerfile for building image of dht-folder, including the node package and main

#Debian system with Go installed

FROM golang

RUN go get -d github.com/fredriklind/D7024E-Lab/dht
WORKDIR $GOPATH/src/github.com/fredriklind/D7024E-Lab/dht

RUN mkdir db
RUN mkdir db/replicas

RUN go build


EXPOSE 3000
EXPOSE 4000
EXPOSE 5000
EXPOSE 8080