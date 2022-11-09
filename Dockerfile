# Build stage
FROM golang:1.16-alpine AS build

WORKDIR /app

RUN mkdir -p gokeyval/cmd/server gokeyval/internal/server gokeyval/scripts

WORKDIR /app/gokeyval

COPY cmd/server ./cmd/server/
COPY internal/server ./internal/server/
COPY go.mod ./
COPY go.sum ./
COPY Makefile ./
COPY ./scripts ./scripts/

RUN apk add --update gcc g++ make && go build -o build/server ./cmd/server

# Deploy stage
FROM alpine

WORKDIR /app

RUN mkdir tls
COPY --from=build /app/gokeyval/build/server ./

RUN chown -R 10010:10010 .

ENV CRL_PATH="tls/list.crl" \ 
    SERVER_KEY="tls/server.key" \ 
    SERVER_CERT="tls/server.crt" \ 
    ROOTCA_CERT="tls/rootCA.crt"

USER 10010

VOLUME /app/tls

EXPOSE 6842

ENTRYPOINT ["./server"]
CMD ["-s", "hash"]