all:
	go build -o ah-authorization
#$(SRC)

run:
	./ah-authorization

#go run $(SRC)

dependencies:
	go get github.com/BurntSushi/toml
	go get -u github.com/gorilla/mux
	go get github.com/go-sql-driver/mysql
	
ahtest:
	curl --insecure -X GET https://127.0.0.1:8445/authorization/echo
	#curl --insecure -X POST -d @ServiceQueryForm.json https://127.0.0.1:4443/serviceregistry/query 


pkcs2pem:
	openssl pkcs12 -in authorization.p12 -nocerts -nodes | sed -ne '/-BEGIN PRIVATE KEY-/,/-END PRIVATE KEY-/p' > authorization.key
	openssl pkcs12 -in authorization.p12 -clcerts -nokeys | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > authorization.pem
	openssl pkcs12 -in authorization.p12 -cacerts -nokeys -chain | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > ca.pem


certs2:
	# ECDSA recommendation key ≥ secp384r1
	# List ECDSA the supported curves (openssl ecparam -list_curves)
	openssl req -x509 -nodes -newkey ec:secp384r1 -keyout server.ecdsa.key -out server.ecdsa.crt -days 3650
	# openssl req -x509 -nodes -newkey ec:<(openssl ecparam -name secp384r1) -keyout server.ecdsa.key -out server.ecdsa.crt -days 3650
	# -pkeyopt ec_paramgen_curve:… / ec:<(openssl ecparam -name …) / -newkey ec:…
	ln -sf server.ecdsa.key server.key
	ln -sf server.ecdsa.crt server.crt

certs:
	# Key considerations for algorithm "RSA" ≥ 2048-bit
	rm -rf server.key server.crt
	openssl genrsa -out server.key 2048

	# Key considerations for algorithm "ECDSA" ≥ secp384r1
	#	List ECDSA the supported curves (openssl ecparam -list_curves)
	openssl ecparam -genkey -name secp384r1 -out server.key
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650

newcerts:
	openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -out cert.pem -keyout key.pem -subj "/C=US/ST=California/L=Mountain View/O=Your Organization/OU=Your Unit/CN=localhost"

