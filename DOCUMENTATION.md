<h1 align="center">
    Festivals Gateway Documentation
</h1>

<p align="center">
  <a href="#overview">Overview</a> •
  <a href="#gateway-route">Gateway-Route</a> •
  <a href="#discovery-route">Discovery-Route</a> •
  <a href="#festivalsapi-route">FestivalsAPI-Route</a> •
  <a href="#database-route">Database-Route</a> •
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
Authorization: Bearer <Header>.<Payload>.<Signatur>
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

[Database-Route](#database-route)

* GET, POST, PATCH, DELETE      `/*`

[FestivalsFilesAPI-Route](#festivalsfilesapi-route)

* GET, POST, PATCH, DELETE      `/*`

------------------------------------------------------------------------------------

## Gateway Status

The **gateway routes** serve status-related information and are available at:

```text
https://gateway.hostname
```

> 📌 Example: `https://gateway.festivalsapp.home`

It is commonly used for health checks, CI/CD diagnostics, or runtime introspection. This route uses
a `server-info` object containing metadata about the currently running binary, such as build time,
Git reference, service name, and version.

### `server-info` object

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

#### GET `/info`

Returns the `server-info`.

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

**Response**

* Codes `200`/`40x`/`50x`
* `data` or `error` field

#### GET `/version`

Returns the release version of the server running.

> In production builds this will have the format `v1.0.2` but
for manual builds this will may be `development`.

**Authorization**
Requires a valid `JWT` token with the user role set to `ADMIN`.

Example:  
  `GET https://gateway.festivalsapp.dev/version`

**Response**

* Codes `200`/`40x`/`50x`
* Server version as a string `text/plain`

#### POST `/update`

Updates to the newest release on github and restarts the service.

Authorization: JWT
  
Example:  
  `POST https://gateway.festivalsapp.dev/update`

Returns
      * The version of the server application.
      * Codes `202`/`40x`/`50x`
      * server version as a string `text/plain`

------------------------------------------------------------------------------------

#### GET `/health`

Authorization: JWT

Example:  
  `GET https://gateway.festivalsapp.dev/health`

Returns
      * Always returns HTTP status code 200
      * Code `200`
      * empty `text/plain`

------------------------------------------------------------------------------------

#### GET `/log`

Returns the service log.

Authorization: JWT

Example:  
  `GET https://gateway.festivalsapp.dev/log`

Returns
      * Returns a string
      * Codes `200`/`40x`/`50x`
      * empty or `text/plain`

------------------------------------------------------------------------------------

#### GET `/log/trace`

Returns the service trace log.

Authorization: JWT

Example:  
  `GET https://gateway.festivalsapp.dev/log/trace`

Returns
      * Returns a string
      * Codes `200`/`40x`/`50x`
      * empty or `text/plain`

------------------------------------------------------------------------------------

## Discovery-Route

The discovery route listens on requests to 'https://discovery.hostname'.

------------------------------------------------------------------------------------

#### POST `/loversear`

Authorization: Service key

Example:  
  `POST https://discovery.festivalsapp.dev/loversear`

Returns
      * Returns nothing on success but a 202 status code.
      * Code `202`/`400`
      * Empty text or `error` field

------------------------------------------------------------------------------------

#### GET `/services`

Authorization: Service key

Example:  
  `GET https://discovery.festivalsapp.dev/services`

Returns
      * Returns the currently available `MonitorNode`s.
      * Code `202`/`40x`
      * `data` or `error` field

------------------------------------------------------------------------------------

## FestivalsAPI-Route

The FestivalsAPI route loadbalances and proxys requests from 'https://api.hostname' to the apropriate services.

####  GET, POST, PATCH, DELETE `/*`

------------------------------------------------------------------------------------

## Database-Route

The database route loadbalances and proxys requests from 'https://database.hostname' to the apropriate services.

####  GET, POST, PATCH, DELETE `/*`

------------------------------------------------------------------------------------

## FestivalsFilesAPI-Route

The FestivalsFilesAPI route loadbalances and proxys requests from 'https://files.hostname' to the apropriate services.

####  GET, POST, PATCH, DELETE `/*`

------------------------------------------------------------------------------------