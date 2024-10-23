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

The service gateway for the FestivalsApp backend, providing access to the [FestivalsAPI](https://github.com/Festivals-App/festivals-server), [static file server](https://github.com/Festivals-App/festivals-fileserver), [database](https://github.com/Festivals-App/festivals-database) and the [website node](https://github.com/Festivals-App/festivals-identity-server) acting as a combined loadbalancer, ingress server and discovery service.

![Figure 1: Architecture Overview Highlighted](https://github.com/Festivals-App/festivals-documentation/blob/main/images/architecture/architecture_overview_gateway.svg "Figure 1: Architecture Overview Highlighted")

<hr />
<p align="center">
  <a href="#development">Development</a> •
  <a href="#deployment">Deployment</a> •
  <a href="#usage">Usage</a> •
  <a href="#documentation">Documentation</a> •
  <a href="#engage">Engage</a> •
  <a href="#licensing">Licensing</a>
</p>
<hr />

## Development

TBA

### Requirements

TBA

### Setup

TBA

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

The install and update scripts should work with any system that uses *systemd* and *firewalld*.

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

### Server

All of the scripts require Ubuntu 20 LTS as the operating system and that the server has already been initialised, see the steps to do that [here](https://github.com/Festivals-App/festivals-documentation/tree/master/deployment/general-vm-setup).

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

Copyright (c) 2020-2023 Simon Gaus.

Licensed under the **GNU Lesser General Public License v3.0** (the "License"); you may not use this file except in compliance with the License.

You may obtain a copy of the License at https://www.gnu.org/licenses/lgpl-3.0.html.

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the [LICENSE](./LICENSE) for the specific language governing permissions and limitations under the License.