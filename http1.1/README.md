# HTTP/1.1 Showcase

## Overview

This project demonstrates the key features of HTTP/1.1 by implementing a server in Go. HTTP/1.1 introduced significant improvements over HTTP/1.0, such as persistent connections, chunked transfer encoding, and request pipelining.

## Features Demonstrated

### Persistent Connections

Enabled by default in HTTP/1.1. The server keeps the connection alive for multiple requests/responses.
Log connections with their state transitions (New, Idle, Active, etc.).

### Chunked Transfer Encoding

Implemented in the /chunked endpoint.
Demonstrates sending data in chunks without knowing the total size in advance.

### Request Pipelining

Simulated in the /pipelining endpoint.
Allows multiple requests to be sent on the same connection without waiting for prior responses.

### Host Header Validation

Ensures the Host header is included in requests, as mandated by HTTP/1.1.

## Running the Server

- Run the server.

```bash
go run main.go
```

Access the server at <http://localhost:8080>.

## Endpoints

| Endpoint    | Description                                                                       |
|-------------|-----------------------------------------------------------------------------------|
| /           | Root endpoint. Validates the Host header and demonstrates persistent connections. |
| /chunked    | Sends a chunked response with simulated delays.                                   |
| /pipelining | Simulates processing delays to demonstrate request pipelining.                    |

## Testing Persistent Connections

Use curl or a browser to send multiple requests,

```bash
curl -v http://localhost:8080/
curl -v http://localhost:8080/chunked
```

Sample request,

```plaintext
➜  http1.1 curl -v http://localhost:8080/
curl -v http://localhost:8080/chunked
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET / HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.4.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Date: Sat, 18 Jan 2025 17:59:35 GMT
< Content-Length: 32
< Content-Type: text/plain; charset=utf-8
< 
Welcome to the HTTP/1.1 Server!
* Connection #0 to host localhost left intact
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /chunked HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.4.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Date: Sat, 18 Jan 2025 17:59:38 GMT
< Transfer-Encoding: chunked
< 
e
Chunk 1: Hello
11
Chunk 2: HTTP/1.1
19
Chunk 3: Chunked response
0

* Connection #0 to host localhost left intact
```

Samples responses,

```plaintext
2025/01/18 22:59:35 New connection established. Active connections: 1
2025/01/18 22:59:35 Connection is active. Active connections: 1
2025/01/18 22:59:35 Received request: GET / from [::1]:50062
2025/01/18 22:59:35 Connection is idle. Active connections: 1
2025/01/18 22:59:35 Connection closed. Active connections: 0
2025/01/18 22:59:35 New connection established. Active connections: 1
2025/01/18 22:59:35 Connection is active. Active connections: 1
2025/01/18 22:59:35 Received request: GET /chunked from [::1]:50063
2025/01/18 22:59:38 Connection is idle. Active connections: 1
2025/01/18 22:59:38 Connection closed. Active connections: 0
```

## Testing Chunked Transfer Encoding

Access the `/chunked` endpoint:

```bash
curl -v http://localhost:8080/chunked
```

Sample request,

```plaintext
➜  http1.1 curl -v http://localhost:8080/chunked
*   Trying [::1]:8080...
* Connected to localhost (::1) port 8080
> GET /chunked HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.4.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Date: Sat, 18 Jan 2025 17:58:51 GMT
< Transfer-Encoding: chunked
< 
e
Chunk 1: Hello
11
Chunk 2: HTTP/1.1
19
Chunk 3: Chunked response
0

* Connection #0 to host localhost left intact
```

Sample response,

```plaintext
2025/01/18 22:58:48 New connection established. Active connections: 1
2025/01/18 22:58:48 Connection is active. Active connections: 1
2025/01/18 22:58:48 Received request: GET /chunked from [::1]:49614
2025/01/18 22:58:51 Connection is idle. Active connections: 1
2025/01/18 22:58:51 Connection closed. Active connections: 0
```

## Testing Request Pipelining

Use netcat to test pipelining:

```bash
(echo -e "GET /pipelining HTTP/1.1\r\nHost: localhost:8080\r\n\r\nGET / HTTP/1.1\r\nHost: localhost:8080\r\n\r\n" | nc localhost 8080)
```

Sample logs,

```plaintext
➜  http1.1 go run main.go
HTTP/1.1 Server is running on port 8080
2025/01/18 22:57:46 New connection established. Active connections: 1
2025/01/18 22:57:46 Connection is active. Active connections: 1
2025/01/18 22:57:46 Received request: GET /pipelining from [::1]:65420
2025/01/18 22:57:48 Connection is idle. Active connections: 1
2025/01/18 22:57:48 Connection is active. Active connections: 1
2025/01/18 22:57:48 Received request: GET / from [::1]:65420
2025/01/18 22:57:48 Connection is idle. Active connections: 1
2025/01/18 22:57:48 Connection closed. Active connections: 0
```

## Logs Example

Logs show connection states and received requests:

```plaintext
New connection established. Active connections: 1
Received request: GET / from 127.0.0.1:54321
Connection is idle. Active connections: 1
Received request: GET /chunked from 127.0.0.1:54321
Connection closed. Active connections: 0
```

### Graceful Shutdown

Press Ctrl+C to stop the server gracefully:

```plaintext
Shutting down server gracefully...
```
