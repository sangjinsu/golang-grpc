SERVER_CN=localhost

# Generate Certificate Authority + Trust Certificate
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "//CN=${SERVER_CN}"

# Generate the Server Private Key server.key
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

# Get a certificate signing request from the CA server.csr
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "//CN=${SERVER_CN}"

# Sign the certificate with the CA we created
openssl x509 -req -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

# Convert the server certificate to .pem format
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem
