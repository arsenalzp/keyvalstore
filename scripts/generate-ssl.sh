#!/bin/bash

CMD=/usr/bin/openssl
TLS_PATH=$1
OPENSSL_CONF=${TLS_PATH}/openssl.conf
EXTENSION_CONF=${TLS_PATH}/extension.conf
CRL=${TLS_PATH}/list.crl
CA_SUBJ="/C=UA/ST=Zaporizhzhia/L=Zaporizhzhia/O=GOKEYVAL/OU=GOKEYVAL/CN=example.com"
CA_CRT=${TLS_PATH}/rootCA.crt
CA_KEY=${TLS_PATH}/rootCA.key
CA_SRL=${TLS_PATH}/rootCA.srl
INDEX=${TLS_PATH}/index.txt
SERVER_CRT=${TLS_PATH}/server.crt
SERVER_KEY=${TLS_PATH}/server.key
SERVER_CSR=${TLS_PATH}/server.csr
CLIENT_CRT=${TLS_PATH}/client.crt
CLIENT_KEY=${TLS_PATH}/client.key
CLIENT_CSR=${TLS_PATH}/client.csr
ENDPOINT_SUBJ="/C=UA/ST=Zaporizhzhia/L=Zaporizhzhia/O=GOKEYVAL/OU=GOKEYVAL/CN=*.example.com"

if ! [[ -z "$(ls -A ${TLS_PATH})" ]]; then
   echo "Removing old keys and certs"
   rm -f ${SERVER_CRT} && \
   rm -f ${SERVER_KEY} && \
   rm -f ${SERVER_CSR} && \
   rm -f ${CA_CRT} && \
   rm -f ${CA_KEY} && \
   rm -f ${CA_SRL} && \
   rm -f ${CLIENT_CRT} && \
   rm -f ${CLIENT_KEY} && \
   rm -f ${CLIENT_CSR} && \
   rm -f ${EXTENSION_CONF} && \
   rm -f ${OPENSSL_CONF}
fi

echo "Generating CA"
${CMD} req -x509 -sha256 -days 3650 -newkey rsa:2048 -nodes -keyout ${CA_KEY} -out ${CA_CRT} -subj ${CA_SUBJ} && \

echo "Preparing openssl.conf file"
cat <<EOF> ${OPENSSL_CONF}
[ca]
default_ca = CA_default
[CA_default]
database = ${TLS_PATH}/index.txt
EOF

echo "Preparing extension file"
cat <<EOF> ${EXTENSION_CONF}
subjectAltName = @san
[san]
DNS.1 = server.example.com
DNS.2 = *.example.com
DNS.3 = localhost
EOF

echo "Preparing index file"
touch ${INDEX}

echo "Generating CRL"
${CMD} ca -gencrl -config ${OPENSSL_CONF} -crldays 3650 -md md5 -out ${CRL} -cert ${CA_CRT} -keyfile ${CA_KEY} && \

echo "Generating server's CSR"
${CMD} req -newkey rsa:2048 -nodes -keyout ${SERVER_KEY} -out ${SERVER_CSR} -subj ${ENDPOINT_SUBJ} && \

echo "Signing server's CSR"
${CMD} x509 -req -CA ${CA_CRT} -CAkey ${CA_KEY} -in ${SERVER_CSR} -out ${SERVER_CRT} -days 365 -CAcreateserial -extfile ${EXTENSION_CONF}

if [[ $2 -eq "cli" ]]; then
  echo "Generating client's CSR"
  ${CMD} req -newkey rsa:2048 -nodes -keyout ${CLIENT_KEY} -out ${CLIENT_CSR} -subj ${ENDPOINT_SUBJ} && \

  echo "Signing client's CSR"
  ${CMD} x509 -req -CA ${CA_CRT} -CAkey ${CA_KEY} -in ${CLIENT_CSR} -out ${CLIENT_CRT} -days 365 -CAcreateserial -extfile ${EXTENSION_CONF}
fi
