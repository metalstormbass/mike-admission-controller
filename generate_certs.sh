# Certificate authority key and certificate
openssl req -new -x509 -days 36500 -nodes -subj '/CN=Mike Admission Controller Webhook' -keyout ca.key -out ca.crt 

# Server key
openssl genrsa -out server.key 2048

# Certificate signing request
openssl req -new -key server.key -subj '/CN=mike-admission-controller.mike-admission-controller.svc' -out server.csr

# Server certificate
bash -c 'openssl x509 -req -days 36500 -extfile <(printf "subjectAltName=DNS:mike-admission-controller.mike-admission-controller.svc") -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt' 

