#!/usr/bin/make --no-print-directory --jobs=1 --environment-overrides -f

CORELIBS_MK_VERSION := v0.1.1

SHELL = /bin/bash
LOCAL_CORELIBS_PATH ?= ..

.PHONY: help version
.PHONY: local unlocal be-update tidy
.PHONY: build clean
.PHONY: test coverage goconvey

#
#: Custom functions
#

define __list_corelibs
$(shell grep -h -v '^module' go.mod \
		| grep 'go-curses/corelibs/' \
		| perl -pe 's!^.*(github\.com/go-curses/corelibs/\S+).*$$!$${1}!' \
		| sort -u -V)
endef

#
#: Actual targets
#

help:
	@echo "# usage: make <help|version>"
	@echo "#        make <local|unlocal|be-update|tidy>"
	@echo "#        make <build|clean>"
	@echo "#        make <test|coverage|goconvey>"

version:
	@echo "# Go-CoreLibs Makefile: ${CORELIBS_MK_VERSION}"

local: export FOUND_LIBS=$(call __list_corelibs)
local:
	@if [ -n "$${FOUND_LIBS}" ]; then \
		for found_lib in $${FOUND_LIBS}; do \
			name=`basename $${found_lib}`; \
			echo "# go mod local ${LOCAL_CORELIBS_PATH}/$${name}"; \
			go mod edit -replace=$${found_lib}=${LOCAL_CORELIBS_PATH}/$${name}; \
		done; \
	else \
		echo "# nothing to do"; \
	fi

unlocal: export FOUND_LIBS=$(call __list_corelibs)
unlocal:
	@if [ -n "$${FOUND_LIBS}" ]; then \
		for found_lib in $${FOUND_LIBS}; do \
			name=`basename $${found_lib}`; \
			echo "# go mod unlocal go-corelibs/$${name}"; \
			go mod edit -dropreplace=$${found_lib}; \
		done; \
	else \
		echo "# nothing to do"; \
	fi

be-update: export GOPROXY=direct
be-update: export FOUND_LIBS=$(call __list_corelibs)
be-update:
	@if [ -n "$${FOUND_LIBS}" ]; then \
		for found_lib in $${FOUND_LIBS}; do \
			name=`basename $${found_lib}`; \
			echo "# go get go-corelibs/$${name}"; \
			go get $${found_lib}@latest; \
		done; \
	else \
		echo "# nothing to do"; \
	fi

tidy:
	@go mod tidy

build:
	@go build -v ./...

clean:
	@rm -fv coverage.{out,html}

test:
	@go test -v ./...

coverage: XDG_OPEN_BIN=`which xdg-open`
coverage: OSX_OPEN_BIN=`which open`
coverage:
	@go test -coverpkg=./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o=coverage.html
	@if [ -n "${XDG_OPEN_BIN}" ]; then \
		${XDG_OPEN_BIN} ./coverage.html; \
	elif [ -n "${OSX_OPEN_BIN}" ]; then \
		${OSX_OPEN_BIN} ./coverage.html; \
	fi

goconvey:
	@echo "# running goconvey... (press <CTRL+c> to stop)"
	@goconvey 2>&1 > /dev/null
