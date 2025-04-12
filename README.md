# Article: RFC 9754 and its impact on API governance

## Project Structure

```bash
├── build
│   ├── krakend
│   │   ├── config
│   │   │   ├── partials          		# Partial configuration files for Krakend
│   │   │   ├── settings          		# General settings for Krakend
│   │   │   │   └── operations.json 	# JSON file defining API operations
│   │   │   └── templates         		# Templates for generating Krakend configurations
│   │   │       └── endpoints.tmpl 		# Template for API endpoints
│   │   ├── Dockerfile            		# Dockerfile to build the Krakend service
│   │   ├── krakend.tmpl          		# Main Krakend configuration template
│   │   └── plugin
│   │       └── deprecated-headers.lua 	# Lua plugin for handling deprecated headers
│   └── parser
│       ├── go.mod                		# Go module definition for the parser
│       ├── go.sum                		# Go dependencies checksum
│       ├── main.go               		# Main entry point for the parser
│       └── README.md             		# Documentation for the parser
├── docker-compose.yaml           		# Docker Compose configuration
├── openapi.yaml                  		# OpenAPI specification for the API
└── README.md                     		# Project documentation
```

- **build/krakend**: Contains all files related to the Krakend API Gateway, including configuration, templates, and plugins.
- **build/parser**: Contains the Go-based parser for processing API specifications.
- **docker-compose.yaml**: Defines the services, networks, and volumes for running the project with Docker Compose.
- **openapi.yaml**: The OpenAPI specification file describing the API endpoints and operations.
- **README.md**: Documentation for the project, including usage instructions and structure.

## Running the Project with Docker Compose

To run the project using Docker Compose, follow these steps:

1. Ensure you have Docker and Docker Compose installed on your system.
3. Start the services using Docker Compose (rebuilding the project every time):
   ```bash
   docker-compose up --build
   ```
4. The services will start, and you can access them as defined in the `docker-compose.yaml` file.

To stop the services, press `Ctrl+C` and run:
```bash
docker-compose down
```
## Making requests to the API

### NOT deprecated operation example

Request
```bash
curl -i -X GET localhost:8080/users/1
```
Response Headers
```bash
HTTP/1.1 200 OK
Content-Length: 250
Content-Type: application/json; charset=utf-8
Date: Sat, 12 Apr 2025 23:23:43 GMT
X-Krakend: Version 2.9.3
X-Krakend-Completed: false
```

DEPRECATED operation example
```bash
curl -i -X DELETE localhost:8080/users/1
```

Response Headers
```bash
HTTP/1.1 200 OK
Content-Length: 253
Content-Type: application/json; charset=utf-8
Date: Sat, 12 Apr 2025 23:24:27 GMT
X-Deprecated-At: @1751327999
X-Deprecated-Link: https://api.example.com/deprecation-policy
X-Deprecated-Sunset: Thu, 01 Jan 2026 00:00:00 UTC
X-Krakend: Version 2.9.3
X-Krakend-Completed: false
```