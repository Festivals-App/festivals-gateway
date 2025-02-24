<!--suppress ALL -->

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

To access the gateway you need to either provide a service key via a custom header or a JWT with your requests authorization header, requests to the loadbalanced services don't need any means of authentication:

```ini
Service-Key:<service-key>
Authorization: Bearer <jwt>
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

## Gateway-Route

The gateway route listens on requests to `https://gateway.hostname` (for Example: <https://gateway.festivalsapp.home>).

## Server Status

Determine the state of the server.
Info object

```json
{
    "BuildTime":      string,
    "GitRef":         string,
    "Version":        string
}
```

------------------------------------------------------------------------------------
#### GET `/info`

* Authorization: JWT

* Example:  
  `GET https://gateway.festivalsapp.dev/info`

* Returns:
      * Returns the info object
      * Codes `200`/`40x`/`50x`
      * `data` or `error` field

------------------------------------------------------------------------------------
#### GET `/version`

* Authorization: JWT

* Example:  
  `GET https://gateway.festivalsapp.dev/version`
 
* Returns:
      * The version of the server application.
      * Codes `200`/`40x`/`50x`
      * server version as a string `text/plain`

------------------------------------------------------------------------------------
#### POST `/update`

Updates to the newest release on github and restarts the service.

 * Authorization: JWT
  
 * Example:  
  `POST https://gateway.festivalsapp.dev/update`

 * Returns
      * The version of the server application.
      * Codes `202`/`40x`/`50x`
      * server version as a string `text/plain`

------------------------------------------------------------------------------------
#### GET `/health`

 * Authorization: JWT
 
 * Example:  
  `GET https://gateway.festivalsapp.dev/health`

 * Returns
      * Always returns HTTP status code 200
      * Code `200`
      * empty `text/plain`

------------------------------------------------------------------------------------
#### GET `/log`

Returns the service log.

 * Authorization: JWT
 
 * Example:  
  `GET https://gateway.festivalsapp.dev/log`

 * Returns
      * Returns a string
      * Codes `200`/`40x`/`50x`
      * empty or `text/plain`

------------------------------------------------------------------------------------
#### GET `/log/trace`

Returns the service trace log.

 * Authorization: JWT
 
 * Example:  
  `GET https://gateway.festivalsapp.dev/log/trace`

 * Returns
      * Returns a string
      * Codes `200`/`40x`/`50x`
      * empty or `text/plain`

------------------------------------------------------------------------------------
## Discovery-Route

The discovery route listens on requests to 'https://discovery.hostname'.

------------------------------------------------------------------------------------
#### POST `/loversear`

 * Authorization: Service key
 
 * Example:  
  `POST https://discovery.festivalsapp.dev/loversear`

 * Returns
      * Returns nothing on success but a 202 status code.
      * Code `202`/`400`
      * Empty text or `error` field

------------------------------------------------------------------------------------
#### GET `/services`

 * Authorization: Service key
 
 * Example:  
  `GET https://discovery.festivalsapp.dev/services`

 * Returns
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