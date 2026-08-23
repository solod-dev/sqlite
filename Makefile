.PHONY: bind example help
.DEFAULT_GOAL := help

# Homebrew keeps sqlite keg-only, so its pkg-config file is off the default path.
PKG_CONFIG := PKG_CONFIG_PATH=/opt/homebrew/opt/sqlite/lib/pkgconfig:$(PKG_CONFIG_PATH) pkg-config

INC_DIR := $(shell $(PKG_CONFIG) --variable=includedir sqlite3 2>/dev/null)

ifeq ($(INC_DIR),)
  INC_DIR := /opt/homebrew/opt/sqlite/include
  FLAGS  := -I/opt/homebrew/opt/sqlite/include -L/opt/homebrew/opt/sqlite/lib
else
  FLAGS  := $(shell $(PKG_CONFIG) --cflags-only-I --libs-only-L sqlite3)
endif

CFLAGS ?= $(FLAGS)

help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  bind                - Generate the bindings"
	@echo "  example name=<name> - Build an example program"

bind:
	@sobind -o libsqlite3/extern.go -pkg=libsqlite3 -I $(INC_DIR) -style=cap -strip=sqlite3_ $(INC_DIR)/sqlite3.h

example:
	@CFLAGS="$(CFLAGS)" so build -check=sanitize -o ./build/$(name) ./example/$(name)
