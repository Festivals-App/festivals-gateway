<h1 align="center">
    Festivals Gateway Documentation
</h1>

<p align="center">
  <a href="#overview">Overview</a> •
  <a href="#gateway-route">Gateway-Route</a> •
  <a href="#discovery-route">Discovery-Route</a> •
  <a href="#festivalsapi-route">FestivalsAPI-Route</a> •
  <a href="#festivalsfilesapi-route">FestivalsFilesAPI-Route</a>
</p>

### Used Languages

* Documentation: `Markdown`, `HTML`
* Server Application: `golang`
* Deployment: `bash`

### Authentication & Authorization

To authenticate to the gateway you need to either provide a service key via a custom header or a JWT
with your requests authorization header, **requests to the loadbalanced services only need valid mTLS certificates**.

```text
Service-Key: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
Authorization: Bearer <JWT>
```

If you have the authorization to call the given endpoint is determined by your
[user role](https://github.com/Festivals-App/festivals-identity-server/blob/master/auth/user.go).

#### Making a request with curl

```bash
curl -H "X-Request-ID: <uuid>" -H "Authorization: Bearer <JWT>" --cacert ca.crt --cert client.crt --key client.key https://gateway.festivalsapp.home/info
```

## Overview

[Gateway-Route](#gateway-route)

* GET              `/info`
* GET              `/version`
* POST             `/update`
* GET              `/health`
* GET              `/log`
* GET              `/log/trace`

[Discovery-Route](#discovery-route)

* POST            `/loversear`
* GET             `/services`

[FestivalsAPI-Route](#festivalsapi-route)

* GET, POST, PATCH, DELETE      `/*`

[FestivalsFilesAPI-Route](#festivalsfilesapi-route)

* GET, POST, PATCH, DELETE      `/*`

------------------------------------------------------------------------------------

## Gateway Route

The **gateway routes** serve status-related information and are available at:

```text
https://gateway.hostname
```

> 📌 Example: `https://gateway.festivalsapp.home`

It is commonly used for health checks, CI/CD diagnostics, or runtime introspection. This route uses
a `server-info` object containing metadata about the currently running binary, such as build time,
Git reference, service name, and version.

**`server-info`** object

```json
{
  "BuildTime": "string",
  "GitRef": "string",
  "Service": "string",
  "Version": "string"
}
```

| Field      | Description                                                                 |
|------------|-----------------------------------------------------------------------------|
| `BuildTime` | Timestamp of the binary build. Format: `Sun Apr 13 13:55:44 UTC 2025`       |
| `GitRef`    | Git reference used for the build. Format: `refs/tags/v2.2.0` [🔗 Git Docs](https://git-scm.com/book/en/v2/Git-Internals-Git-References) |
| `Service`   | Service identifier. Matches a defined [Service type](https://github.com/Festivals-App/festivals-server-tools/blob/main/heartbeattools.go) |
| `Version`   | Version tag of the deployed binary. Format: `v2.2.0`                        |

> In production builds, these values are injected at build time and reflect the deployment source and context.

------------------------------------------------------------------------------------

#### GET `/info`

Returns the `server-info`.

Example:  
  `GET https://gateway.festivalsapp.home/info`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* `data` or `error` field
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

#### GET `/version`

Returns the release version of the server.

> In production builds this will have the format `v2.2.0` but
for manual builds this will may be `development`.

Example:  
  `GET https://gateway.festivalsapp.home/version`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Server version as a string `text/plain`.
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

#### POST `/update`

Updates to the newest release on github and restarts the service.

Example:  
  `POST https://gateway.festivalsapp.home/update`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* The current version and the version the server is updated to as a string `text/plain`. Format: `v2.1.3 => v2.2.0`
* Codes `202`/`40x`/`50x`

------------------------------------------------------------------------------------

#### GET `/health`

A simple health check endpoint that returns a `200 OK` status if the service is running and able to respond.

Example:  
  `GET https://gateway.festivalsapp.home/health`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Always returns `200 OK`

------------------------------------------------------------------------------------

#### GET `/log`

Returns the info log file as a string, containing all log messages except trace log entries.
See [loggertools](https://github.com/Festivals-App/festivals-server-tools/blob/main/DOCUMENTATION.md#loggertools) for log format.

Example:  
  `GET https://gateway.festivalsapp.home/log`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Returns a string as `text/plain`
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

#### GET `/log/trace`

Returns the trace log file as a string, containing all remote calls to the server.
See [loggertools](https://github.com/Festivals-App/festivals-server-tools/blob/main/DOCUMENTATION.md#loggertools) for log format.

Example:  
  `GET https://gateway.festivalsapp.home/log/trace`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Returns a string as `text/plain`
* Codes `200`/`40x`/`50x`

------------------------------------------------------------------------------------

## Discovery-Route

The **discovery routes** provide discovery information and are available at:

```text
https://discovery.hostname
```

> 📌 Example: `https://discovery.festivalsapp.home`

The discovery route is commonly used to expose information that helps clients automatically locate
and interact with services in a distributed system. This route uses a `MonitorNode` object containing
metadata about a known service including the type of service, the location and the last time the service
did send a heartbeat.

**`MonitorNode`** object

```json
{
  "Type":    "string",
  "Location": "string",
  "LastSeen": "string"
}
```

| Field      | Description                                                                 |
|------------|-----------------------------------------------------------------------------|
| `Type`     | Service identifier. Matches a defined [Service type](https://github.com/Festivals-App/festivals-server-tools/blob/main/heartbeattools.go)  |
| `Location` | The URL to the service. Format: `https://server-0.festivalsapp.home:10439` |
| `LastSeen` | The last time the service was registered. Format: `2025-04-16T00:16:17.079845416Z` |

------------------------------------------------------------------------------------

#### POST `/loversear`

Registers the heartbeat call from other services.

Example:  
  `POST https://discovery.festivalsapp.home/loversear`

**Authorization**
Requires a valid `Service-Key`.

**Response**

* Returns `202 Accepted` on success and `error` on failure.
* Code `202`/`400`

------------------------------------------------------------------------------------

#### GET `/services`

Returns all known services as a list of `MonitorNode`s.

Example:  
  `GET https://discovery.festivalsapp.home/services`

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* `data` or `error` field
* Code `202`/`40x`

------------------------------------------------------------------------------------

## FestivalsAPI-Route

The **FestivalsAPI route** loadbalances and proxys requests to available festivals-server services.

```text
https://api.hostname
```

> 📌 Example: `https://api.festivalsapp.home/*`

#### GET, POST, PATCH, DELETE `/*`

------------------------------------------------------------------------------------

## FestivalsFilesAPI-Route

The **FestivalsFilesAPI** route loadbalances and proxys requests to available festivals-fileserver services.

```text
https://files.hostname
```

> 📌 Example: `https://files.festivalsapp.home/*`

#### GET, POST, PATCH, DELETE  `/*`

------------------------------------------------------------------------------------
