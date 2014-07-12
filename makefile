#!/usr/bin/make -f

SHELL=/bin/bash
bin=bin
name=push

all: build

build:
	go build -v -o bin/$(name)

release:
	GOARCH=386 OS=darwin go build -o $(bin)/$(name)-Darwin-i386; \
	GOARCH=amd64 OS=darwin go build -o $(bin)/$(name)-Darwin-x86_64; \
	GOARCH=386 OS=linux go build -o $(bin)/$(name)-Linux-i386; \
	GOARCH=amd64 OS=linux go build -o $(bin)/$(name)-Linux-amd64; \
	GOARCH=arm OS=linux go build -o $(bin)/$(name)-Linux-armv6l; \
	GOARCH=386 OS=freebsd go build -o $(bin)/$(name)-FreeBSD-i386; \
	GOARCH=amd64 OS=linux go build -o $(bin)/$(name)-FreeBSD-amd64; \
	GOARCH=386 OS=windws go build -o $(bin)/$(name)-windows-i386.exe
	GOARCH=amd64 OS=windows go build -o $(bin)/$(name)-windows-amd64.exe

clean:
	go clean -x

remove:
	go clean -i