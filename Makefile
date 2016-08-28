SHELL := /bin/bash
SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

all: clean im_client testplugin

clean:
	[[ -e $(SELF_DIR)im_client ]] && rm $(SELF_DIR)im_client || true
	[[ -e $(SELF_DIR)testplugin ]] && rm $(SELF_DIR)testplugin || true

im_client:
	go build $(SELF_DIR)cmd/im_client

testplugin:
	go build $(SELF_DIR)cmd/testplugin


