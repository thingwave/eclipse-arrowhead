
all:
	go build -o ah-serviceregistry
#$(SRC)

run:
	./ah-serviceregistry

clean:
	rm -rf ah-serviceregistry

pkcs2pem:
	openssl pkcs12 -in service_registry.p12 -nocerts -nodes | sed -ne '/-BEGIN PRIVATE KEY-/,/-END PRIVATE KEY-/p' > service_registry.key
	openssl pkcs12 -in service_registry.p12 -clcerts -nokeys | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > service_registry.pem
	openssl pkcs12 -in service_registry.p12 -cacerts -nokeys -chain | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > ca.pem

