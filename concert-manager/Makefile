# Makefile
all: setup buildtui buildserver
tui: setup buildtui
server: setup buildserver

setup:
	mkdir -p ./build

test:

buildserver:
	go build -o ./build/cm-server ./cmd/server

runserver:
	. ./env.sh && ./build/cm-server $(ARGS)

buildtui:
	go build -o ./build/cm-tui ./cmd/tui

runtui:
	. ./env.sh && ./build/cm-tui $(ARGS)
