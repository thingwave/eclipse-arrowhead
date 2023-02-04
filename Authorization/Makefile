
all:
	go build -o ah-authorization

all-arm64:
        env GOOS=linux GOARCH=arm GOARM=5 go build -o ah-authorization

run:
	./ah-authorization

ahtest:
	curl --insecure -X GET https://127.0.0.1:8445/authorization/echo

pkcs2pem:
	openssl pkcs12 -in authorization.p12 -nocerts -nodes | sed -ne '/-BEGIN PRIVATE KEY-/,/-END PRIVATE KEY-/p' > authorization.key
	openssl pkcs12 -in authorization.p12 -clcerts -nokeys | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > authorization.pem
	openssl pkcs12 -in authorization.p12 -cacerts -nokeys -chain | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > ca.pem
