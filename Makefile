CURR_DIR := $(shell pwd)
ARCH := $(shell uname -m)
INSTALL := ./dist/${ARCH}
TARGET := ${INSTALL}/ezbot

ENTRY_EZBOT := ./cmd/ezbot/main.go
BIN_EZBOT := ./bin/ezbot

OPT_LIB_ENV := LD_LIBRARY_PATH=`pwd`/opt/lib/${ARCH}:${LD_LIBRARY_PATH}

help:
	@echo "make [ezbot|all|test|run|air|install|clean]"
	@mkdir -p ./bin

all: ezbot
	@echo "DONE."

ezbot:
	@echo "build ezbot"
	@#go build -o ${BIN_EZBOT}.${ARCH} ${ENTRY_EZBOT}
ifeq ($(ARCH), x86_64)
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ${BIN_EZBOT}.x86_64 ${ENTRY_EZBOT}
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc CGO_LDFLAGS="-L ./opt/lib/aarch64" go build -o ${BIN_EZBOT}.aarch64 ${ENTRY_EZBOT}
else
	@echo "Please run in x86_64 linux OS."
endif

test:
	@${OPT_LIB_ENV} && cd ./internal/strategies && go test -cover -v .

run:
	@${OPT_LIB_ENV} go run ${ENTRY_EZBOT} -config ezbot.yaml

# run and auto-reload, need on virtual env
# go install github.com/cosmtrek/air@latest
air:
	air -- -config ezbot.yaml

install:
	@echo "install ezbot..."
	cp -af ${BIN_EZBOT}.${ARCH} ${TARGET}/bin/ezbot
	cp -af ./opt/lib/${ARCH}/libxxx.so ${TARGET}/lib/
	cp -af ezbot.yaml ${TARGET}/config/

clean:
	@echo "do clean..."
	@rm -rf ./bin/*
