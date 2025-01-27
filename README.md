<p align="center">
   <a href="https://github.com/festivals-app/festivals-gateway/commits/" title="Last Commit"><img src="https://img.shields.io/github/last-commit/festivals-app/festivals-gateway?style=flat"></a>
   <a href="https://github.com/festivals-app/festivals-gateway/issues" title="Open Issues"><img src="https://img.shields.io/github/issues/festivals-app/festivals-gateway?style=flat"></a>
  <a href="https://github.com/festivals-app/festivals-gateway" title="SLSA Level 2"><img src="https://img.shields.io/badge/SLSA-Level_2-blue"></a>
   <a href="./LICENSE" title="License"><img src="https://img.shields.io/github/license/festivals-app/festivals-gateway.svg"></a>
</p>

<h1 align="center">
  <br/><br/>
    FestivalsApp Gateway
  <br/><br/>
</h1>

The service gateway for the FestivalsApp backend, providing access to the [FestivalsAPI](https://github.com/Festivals-App/festivals-server), [FestivalsFilesAPI](https://github.com/Festivals-App/festivals-fileserver), [database](https://github.com/Festivals-App/festivals-database) and the [website node](https://github.com/Festivals-App/festivals-identity-server) acting as a ingress server, loadbalancer and discovery service.

![Figure 1: Architecture Overview Highlighted](https://github.com/Festivals-App/festivals-documentation/blob/main/images/architecture/architecture_overview_gateway.svg "Figure 1: Architecture Overview Highlighted")

<hr />
<p align="center">
  <a href="#development">Development</a> •
  <a href="#deployment">Deployment</a> •
  <a href="#engage">Engage</a>
</p>
<hr />

## Development

TBA

To find out more about the architecture and technical information see the [ARCHITECTURE](./ARCHITECTURE.md) document. The general documentation for the Festivals App is in the [festivals-documentation](https://github.com/festivals-app/festivals-documentation) repository. The documentation repository contains architecture information, general deployment documentation, templates and other helpful documents.

### Requirements

- [Golang](https://go.dev/) Version 1.23.5+
- [Visual Studio Code](https://code.visualstudio.com/download) 1.96.0+
  - Plugin recommendations are managed via [workspace recommendations](https://code.visualstudio.com/docs/editor/extension-marketplace#_recommended-extensions).
- [Bash script](https://en.wikipedia.org/wiki/Bash_(Unix_shell)) friendly environment

## Deployment

The Go binaries are able to run without system dependencies so there are not many requirements for the system to run the festivals-gateway binary.
The config file needs to be placed at `/etc/festivals-gateway.conf` or the template config file needs to be present in the directory the binary runs in.

You also need to provide certificates in the right format and location:

- The default path to the root CA certificate is  `/usr/local/festivals-gateway/ca.crt`
- The default path to the server certificate is   `/usr/local/festivals-gateway/server.crt`
- The default path to the corresponding key is    `/usr/local/festivals-gateway/server.key`

Where the root CA certificate is required to validate incoming requests and the server certificate and key is requires to make outgoing connections.
For instructions on how to manage and create the certificates see the [festivals-pki](https://github.com/Festivals-App/festivals-pki) repository.

### VM deployment

The install and update scripts should work with any system that uses *systemd* and *ufw*.

Installing
```bash
curl -o install.sh https://raw.githubusercontent.com/Festivals-App/festivals-gateway/main/operation/install.sh
chmod +x install.sh
sudo ./install.sh
```
Updating
```bash
curl -o update.sh https://raw.githubusercontent.com/Festivals-App/festivals-gateway/main/operation/update.sh
chmod +x update.sh
sudo ./update.sh
```

#### Build and run using make

```bash
make build
make run
# Default API Endpoint : http://localhost:10439
```

## Usage

TBA

base/health
base/version
base/info
base/log

discovery.base/services
discovery.base/loversear

api.base/*

files.base/*

### Documentation

The gateway is documented in detail [here](./DOCUMENTATION.md).

The full documentation for the Festivals App is in the [festivals-documentation](https://github.com/festivals-app/festivals-documentation) repository. 
The documentation repository contains technical documents, architecture information, UI/UX specifications, and whitepapers related to this implementation.

## Engage

I welcome every contribution, whether it is a pull request or a fixed typo. The best place to discuss questions and suggestions regarding the festivals-gateway is the [issues](https://github.com/festivals-app/festivals-gateway/issues/) section. More general information and a good starting point if you want to get involved is the [festival-documentation](https://github.com/Festivals-App/festivals-documentation) repository.

The following channels are available for discussions, feedback, and support requests:

| Type                     | Channel                                                |
| ------------------------ | ------------------------------------------------------ |
| **General Discussion**   | <a href="https://github.com/festivals-app/festivals-documentation/issues/new/choose" title="General Discussion"><img src="https://img.shields.io/github/issues/festivals-app/festivals-documentation/question.svg?style=flat-square"></a> </a>   |
| **Other Requests**    | <a href="mailto:simon.cay.gaus@gmail.com" title="Email me"><img src="https://img.shields.io/badge/email-Simon-green?logo=mail.ru&style=flat-square&logoColor=white"></a>   |

## Licensing

Copyright (c) 2020-2025 Simon Gaus. Licensed under the [**GNU Lesser General Public License v3.0**](./LICENSE)