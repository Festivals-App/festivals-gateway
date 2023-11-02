<p align="center">
   <a href="https://github.com/festivals-app/festivals-gateway/commits/" title="Last Commit"><img src="https://img.shields.io/github/last-commit/festivals-app/festivals-gateway?style=flat"></a>
   <a href="https://github.com/festivals-app/festivals-gateway/issues" title="Open Issues"><img src="https://img.shields.io/github/issues/festivals-app/festivals-gateway?style=flat"></a>
   <a href="./LICENSE" title="License"><img src="https://img.shields.io/github/license/festivals-app/festivals-gateway.svg"></a>
</p>

<h1 align="center">
  <br/><br/>
    Festivals Gateway Server
  <br/><br/>
</h1>

<p align="center">
  <a href="#development">Development</a> •
  <a href="#deployment">Deployment</a> •
  <a href="#usage">Usage</a> •
  <a href="#documentation">Documentation</a> •
  <a href="#engage">Engage</a> •
  <a href="#licensing">Licensing</a>
</p>

The service gateway for the FestivalsAPI, providing access to the [FestivalsAPI](https://github.com/Festivals-App/festivals-server), [Website](https://github.com/Festivals-App/festivals-website), [static file server](https://github.com/Festivals-App/festivals-fileserver) and the [identity service](https://github.com/Festivals-App/festivals-identity-server).

![Figure 1: Architecture Overview Highlighted](https://github.com/Festivals-App/festivals-documentation/blob/main/images/architecture/overview_gate.png "Figure 1: Architecture Overview Highlighted")

## Development

TBA

### Requirements

TBA

### Setup

TBA

## Deployment

### VM deployment

The install, update and uninstall scripts should work with any system that uses *systemd* and *firewalld*.
Additionally the scripts will somewhat work under macOS but won't configure the firewall or launch service.

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

TBA

### Docker

TBA

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