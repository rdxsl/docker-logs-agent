#!/bin/bash
# taking the idea from Nick Craig-Wood https://gist.github.com/ncw/9253562
# call this script with an email address (valid or not).
# like:
# ./makecert.sh foo@foo.com

if [ "$1" == "" ]; then
    echo "Need email as argument"
    exit 1
fi
EMAIL=$1

CERTPATH="../conf/dev"
mkdir -p $CERTPATH

if [ "$2" != "" ]; then
    echo "Need email as argument"
    CERTPATH=$2
fi

rm -rf certs
mkdir certs
cd certs

echo "make CA"
PRIVKEY="test"
openssl req -new -x509 -days 365 -keyout ca.key -out ca.pem -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=www.random.com/emailAddress=KryptoKings@random.com" -passout pass:$PRIVKEY

echo "make server cert"
openssl req -new -nodes -x509 -out server.pem -keyout server.key -days 3650 -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=www.random.com/emailAddress=${EMAIL}"

echo "make client cert"

openssl genrsa -out client.key 2048
echo "00" > ca.srl
openssl req -sha1 -key client.key -new -out client.req -subj "/C=DE/ST=NRW/L=Earth/O=Random Company/OU=IT/CN=client.com/emailAddress=${EMAIL}"
# Adding -addtrust clientAuth makes certificates Go can't read
openssl x509 -req -days 365 -in client.req -CA ca.pem -CAkey ca.key -passin pass:$PRIVKEY -out client.pem # -addtrust clientAuth

openssl x509 -extfile ../openssl.conf -extensions ssl_client -req -days 365 -in client.req -CA ca.pem -CAkey ca.key -passin pass:$PRIVKEY -out client.pem

cd ..
rm -rf ${CERTPATH}/certs
mv certs ${CERTPATH}
