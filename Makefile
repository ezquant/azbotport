CURR_DIR := $(shell pwd)
ARCH := $(shell uname -m)
INSTALL := ./dist/${ARCH}
TARGET := ${INSTALL}/azbot

ENTRY_azbot := ./cmd/azbot/main.go
BIN_azbot := ./bin/azbot

OPT_LIB_ENV := LD_LIBRARY_PATH=`pwd`/opt/lib/${ARCH}:${LD_LIBRARY_PATH}

help:
	@echo "make [azbot|all|test|run|air|install|clean]"
	@mkdir -p ./bin

all: azbot
	@echo "DONE."

azbot:
	@echo "build azbot"
	@#go build -o ${BIN_azbot}.${ARCH} ${ENTRY_azbot}
ifeq ($(ARCH), x86_64)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BIN_azbot}.x86_64 ${ENTRY_azbot}
	@#CGO_ENABLED=0 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc CGO_LDFLAGS="-L ./opt/lib/aarch64" go build -o ${BIN_azbot}.aarch64 ${ENTRY_azbot}
else
	@echo "Please run in x86_64 linux OS."
endif

test:
	@${OPT_LIB_ENV} && cd ./internal/strategies && go test -cover -v .

run:
	@${OPT_LIB_ENV} go run ${ENTRY_azbot} test -config user_data/config.yml

# run and auto-reload, need on virtual env
# go install github.com/cosmtrek/air@latest
air:
	air -- test -config user_data/config.yml

install:
	@echo "install azbot..."
	cp -af ${BIN_azbot}.${ARCH} ${TARGET}/bin/azbot
	cp -af ./opt/lib/${ARCH}/libxxx.so ${TARGET}/lib/
	cp -af azbot.yaml ${TARGET}/config/

clean:
	@echo "do clean..."
	@rm -rf ./bin/*
