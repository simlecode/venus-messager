SHELL=/usr/bin/env bash

all:
	rm -rf ./venus-message
	go build -o venus-message main.go