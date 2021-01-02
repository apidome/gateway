# APIDome Gateway
A security gateway for cloud web applications.

# Intoduction
### A cross platform API gateway
At its base, APIDome developed for container environments but can also be installed in a bare metal machines and virtual machines. All it requires is the [Golang runtime environment](https://golang.org/) to be installed on the hosting platform.
In its current version, the gateway can:
- Handle SSL or plain text connections
- Request for client certificates
- Validate and Filter HTTP requests payloads based on a json schema


# Getting started

# Configuration
### Structure

```js
{
    // This configuration section determines how the gateway will communicate
    // with the outer world.
    "out": {
        "port": "3000",

        // Boolean. If true, the gateway will listen for https connections.
        "ssl": true,

        // Relative path to a certificate file (relevant only if "ssl" is true).
        "certPath": "",

        // Relative path to a key file (relevant only if "ssl" is true).
        "keyPath": ""
    },
    // This configuration section determines how the gateway will communicate
    // with the entities that it protects.
    "in": {

        // Each item in this array is a an entity that the gateway will proxy requests to.
        "targets": [
            {
                "host": "rapidapi.com",
                "port": "443",
                "ssl": true,

                // Determines whether the gateway should request for a user certicicate for
                // requests that targeted to this entity.
                "clientAuth": false,

                // A list of APIs that the entity serves.
                "apis": [
                    {
                        // Supported API types - for now supports REST APIs only.
                        "type": "REST",

                        // The spec version that the gateway should rely on.
                        "version": "<version>",

                        // A list of endpoints that the API serves.
                        // (GraphQL APIs should have only one endpoint)
                        "endpoints": [
                            {
                                "path": "<some_api_endpoint>",

                                // Any HTTP method (in the future will accept an array of method).
                                "method": "GET",

                                // Path to a schema that tells the gateway how to validate requests.
                                "schema": "schemas/schema1.json"
                            }
                        ],

                        // A set of configuration that configures the API validator behaviour.
                        "validator": {
                            // Boolean. If true, the validator will not block requests, only log.
                            "monitor": true
                        }
                    }
                ]
            }
        ]
    }
}
```
