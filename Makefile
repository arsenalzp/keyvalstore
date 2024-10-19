SERVER_BIN=server
SERVER_SRC=${PWD}/cmd/server
CLIENT_BIN=cli
CLIENT_SRC=${PWD}/cmd/cli
GENERATE_TLS=${PWD}/scripts/generate-ssl.sh

default: server

all: server cli

server:
	mkdir -p ${PWD}/build
	go mod download
	go build -o ${PWD}/build/${SERVER_BIN} ${SERVER_SRC}
	sh -c '${GENERATE_TLS} ${PWD}/build'

cli:
	mkdir -p ${PWD}/build
	go mod download
	go build -o ${PWD}/build/${CLIENT_BIN} ${CLIENT_SRC}
	sh -c '${GENERATE_TLS} ${PWD}/build "cli"'

clean:
	rm -rf ${PWD}/build
