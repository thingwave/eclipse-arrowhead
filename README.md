# Introduction

This is a Go based implementation of the core systems in Eclipse Arrowhead. The primary source for documentation, code and certificates is the [Core-Java-Spring](https://github.com/eclipse-arrowhead/core-java-spring).

The web page is Eclipse Arrowhead can be found [here](https://www.arrowhead.eu).

# Build
To build all core systems, type:
```
    > make all
```

# Certificate management
This implementation relies heavy on PEM certificates. Documentation will be added on how to convert original PKCS#12 certificates used by the Java Spring reference implementation of Eclipse Arrowhead to PEM certificates.

# Todo list

Below is a list if some of the planned features:

* Debian package support
* RHEL package support
* More core systems
* Token support
* ...

# Status:

## ArrowheadClient
Basic HTTPS client with JSON support. Must add support for Serviceregistry and Orchestration

## ServiceRegistry
Implements most features.

## Orchestrator
Implements most features.

## Authorization
Implements most features, excekt intercloud and token management.

## CertificateAuthority
Implements no features.

## DataManager
Implements all features.

