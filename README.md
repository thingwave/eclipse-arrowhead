# Introduction

This is a Go based implementation of the core systems in Eclipse Arrowhead. The primary source for documentation, code and certificates is the [Core-Java-Spring](https://github.com/eclipse-arrowhead/core-java-spring).

The web page is Eclipse Arrowhead can be found [here](https://www.arrowhead.eu).

# Build
To build all core systems, type:
```
    > make all
```

# Certificate management
This implementation relies heavily on PEM certificates. Documentation will be added on how to convert original PKCS#12 certificates used by the Java Spring reference implementation of Eclipse Arrowhead to PEM certificates.

# Eclipse Arrowhead Control app
ThingWave have developed a helper tool called ahctl to interact with a local cloud. More information is available in the GitHub repo for [ahctl](https://github.com/thingwave/ahctl).

# Todo list

Below is a list of some of the planned features:

* Debian package support
* RHEL package support
* More core systems
* Token support
* ...

# Status:

## ServiceRegistry
Implements most features.

## Orchestrator
Implements most features.

## Authorization
Implements most features, except intercloud and token management.

## DataManager
Implements all features.

