##  KEYVALSTORE - Go key-value storage

The simple key-value storage that uses either hash-table or 
SQLite underlying storage.

It uses client-server architecture, a set of commands are sent by using TLS transport.
Node.js and golang clients are provided as well.

PKI is used for end-user authentication with; Certificate Revocation List capability is enabled.

### Message format:
Message size: 771 byte

| Command   | Key       | Value     | EOT
|-----------|-----------|-----------|---------|
| Bytes     | Bytes     | Bytes     | Bytes
| 0-3       | 3-257     | 257-771   | 771

### A set of commands:
+ SET - set a value to a key
+ GET - get a value of a given key
+ DEL - delete a key and its value
+ EXPORT - export all key-value data from the server in JSON format
+ IMPORT - import key-value data to the server in JSON format

### Building KEYVALSTORE

Project can be compiled and used for Linux, BSD, Windows, for both x86_64 and x86

Build CLI and server binary:
```
  make all
```

Build CLI only:
```
  make cli
```

Build server only:
```
  make server
```

Clean all results:
```
  make clean
```

To generate SSL keys, certificates and openssl configs:
```
  ./scripts/generate-ssl.sh path
```

For Container deployment Dockerfile is provided.

### Run a server
To run the server it is necessary to define CA certificate, the server certificate and key, CRL file and underlying hash-table storage:
```
  cd build && \
  CRL_PATH="./list.crl" \
    SERVER_KEY="./server.key" \
    SERVER_CERT="./server.crt" \
    ROOTCA_CERT="./rootCA.crt" \
    SERVICE_STORAGE="hash" \
    ./server
```

for SQLite sotrage, set SERVICE_STORAGE to "sqlite :
```
  SERVICE_STORAGE="sqlite" \
    CRL_PATH="./list.crl" \
    SERVER_KEY="./server.key" \
    SERVER_CERT="./server.crt" \
    ROOTCA_CERT="./rootCA.crt" \
    ./server
```

If the server need to be bound to the certain NIC, define SERVICE_NIC env:
```
  SERVICE_NIC="lo" \
    CRL_PATH="./list.crl" \
    SERVER_KEY="./server.key" \
    SERVER_CERT="./server.crt" \
    ROOTCA_CERT="./rootCA.crt" \
    SERVICE_STORAGE="hash" \
    ./server
```

To let the server lstening on the specific port, define SERVICE_PORT env:
```
  SERVICE_PORT="1234" \
    CRL_PATH="./list.crl" \
    SERVER_KEY="./server.key" \
    SERVER_CERT="./server.crt" \
    ROOTCA_CERT="./rootCA.crt" \
    SERVICE_STORAGE="hash" \
    ./server
```

### Use CLI
Communication with the serer is done by CLI.
The following parameters are required:
+ client certificate
+ client private key
+ CA certificate
+ server address and port in format server:port
```
  cd build && \
  ./cli --cert ./client.crt --key ./client.key --CAcert ./rootCA.crt -s server:port COMMAND
```

SET command:
```
  ./cli --cert ./client.crt --key ./client.key --CAcert ./rootCA.crt -s server:port set key value
```

GET command:
```
  ./cli --cert ./client.crt --key ./client.key --CAcert ./rootCA.crt -s server:port get key
```

DEL command:
```
  ./cli --cert ./client.crt --key ./client.key --CAcert ./rootCA.crt -s server:port del key
```

IMPORT:
```
  echo '[{"key":"key1","value":"val1UPD"}]' |./cli --cert ./client.crt --key ./client.key \
    --CAcert ./rootCA.crt -s 127.0.0.1:6842 import
```

EXPORT:
```
  ./cli --cert ./client.crt --key ./client.key --CAcert ./rootCA.crt -s 127.0.0.1:6842 export
```

### TODO
+ add expiration capability 
### DISCLAIMER
Code is provided AS IS under BSD license

No code tests are provided!
