
all:
	go build -o ah-certificateauthority

run:
	./ah-certificateauthority

clean:
	rm -rf ah-certificateauthority

pkcs2pem:
	openssl pkcs12 -in certificateauthority.p12 -nocerts -nodes | sed -ne '/-BEGIN PRIVATE KEY-/,/-END PRIVATE KEY-/p' > certificateauthority.key
	openssl pkcs12 -in certificateauthority.p12 -clcerts -nokeys | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > certificateauthority.pem
	openssl pkcs12 -in certificateauthority.p12 -cacerts -nokeys -chain | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > ca.pem
