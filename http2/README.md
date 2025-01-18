# HTTP2 Showcase

## Features Highlighted in the HTTP/2 Server

1. **Multiplexed Streams**: Multiple messages are sent over a single connection, demonstrating the elimination of head-of-line blocking from HTTP/1.1.
2. **Header Compression (HPACK)**: Multiple headers are sent to demonstrate HPACK's ability to reduce header size.
3. **Stream Prioritization**: URL query parameter `priority=high` illustrates prioritized responses.
4. **Server Push**: Demonstrates pushing related resources (`style.css` and `script.js`) to the client.

## Endpoints

1. `/multiplex`: Showcases multiplexed responses.
2. `/headers`: Adds multiple headers to observe HPACK.
3. `/priority?priority=high`: High-priority responses.
4. `/push`: Pushes static resources (`style.css` and `script.js`).

## Prequisites

Before testing the server, ensure you have the following,

- Golang installed (1.16 or higher recommended).
- Generate TLS certificate and private key by using the below command,

```bash
openssl req -x509 -newkey rsa:2048 -nodes -keyout server.key -out server.crt -days 365 -subj "/CN=localhost"
```

- **HTTP/2-compatible tools**,
  - Use curl (version 7.47.0 or newer) or `h2load` for testing.
  - Modern web browsers like Chrome or Firefox also support HTTP/2.

## Running the Server

Start the server by running,

```sh
go run main.go
```

## Testing the Features

### 3.1. Multiplexed Streams

- **What to expect**: Responses for multiple requests are sent simultaneously over a single connection.
- **How to test**:
  - Open multiple tabs in your browser and navigate to `https://localhost:8443/multiplex`.
  - Using curl, `curl -k --http2 https://localhost:8443/multiplex`
  - Run multiple commands simultaneously to observe multiplexing in action.

Sample request,

```plaintext
➜  http2 git:(main) ✗ curl -k --http2 https://localhost:8443/multiplex
```

Sample response,

```plaintext
Message 1 from streams
Message 2 from streams
Message 3 from streams
Message 4 from streams
Message 5 from streams
```

### 3.2. Header Compression (HPACK)

- **What to expect**: The headers sent by the server are compressed, reducing overhead.
- **How to test**:
  - Use curl to view headers, `curl -k --http2 -v https://localhost:8443/headers`.
  - Observe the smaller header sizes (requires a tool like Wireshark for detailed packet inspection).

Sample request,

```plaintext
➜  http2 git:(main) ✗ curl -k --http2 -v https://localhost:8443/headers
```

Sample response,

```plaintext
*   Trying [::1]:8443...
* Connected to localhost (::1) port 8443
* ALPN: curl offers h2,http/1.1
* (304) (OUT), TLS handshake, Client hello (1):
* (304) (IN), TLS handshake, Server hello (2):
* (304) (IN), TLS handshake, Unknown (8):
* (304) (IN), TLS handshake, Certificate (11):
* (304) (IN), TLS handshake, CERT verify (15):
* (304) (IN), TLS handshake, Finished (20):
* (304) (OUT), TLS handshake, Finished (20):
* SSL connection using TLSv1.3 / AEAD-CHACHA20-POLY1305-SHA256
* ALPN: server accepted h2
* Server certificate:
*  subject: CN=localhost
*  start date: Jan 18 18:26:53 2025 GMT
*  expire date: Jan 18 18:26:53 2026 GMT
*  issuer: CN=localhost
*  SSL certificate verify result: self signed certificate (18), continuing anyway.
* using HTTP/2
* [HTTP/2] [1] OPENED stream for https://localhost:8443/headers
* [HTTP/2] [1] [:method: GET]
* [HTTP/2] [1] [:scheme: https]
* [HTTP/2] [1] [:authority: localhost:8443]
* [HTTP/2] [1] [:path: /headers]
* [HTTP/2] [1] [user-agent: curl/8.4.0]
* [HTTP/2] [1] [accept: */*]
> GET /headers HTTP/2
> Host: localhost:8443
> User-Agent: curl/8.4.0
> Accept: */*
> 
< HTTP/2 200 
< x-header-1: Value 1
< x-header-10: Value 10
< x-header-2: Value 2
< x-header-3: Value 3
< x-header-4: Value 4
< x-header-5: Value 5
< x-header-6: Value 6
< x-header-7: Value 7
< x-header-8: Value 8
< x-header-9: Value 9
< content-type: text/plain; charset=utf-8
< content-length: 44
< date: Sat, 18 Jan 2025 18:28:54 GMT
< 
Headers added. Check with an HTTP/2 client.
* Connection #0 to host localhost left intact
```

### 3.3. Stream Prioritization

- **What to expect**: High-priority requests are processed before lower-priority ones.
- How to test:
  - Test low-priority requests, `curl -k --http2 https://localhost:8443/priority`
  - Test high-priority requests, `curl -k --http2 "https://localhost:8443/priority?priority=high"`
  - Observe how high-priority requests are responded to faster.

Sample request,

```plaintext
➜  http2 git:(main) ✗  curl -k --http2 https://localhost:8443/priority
```

Sample response,

```plaintext
Default priority stream # takes 3 seconds to respond
Stream priority demonstration. Check with an HTTP/2 client.
```

Sample request for high-priority,

```plaintext
➜  http2 git:(main) ✗ curl -k --http2 "https://localhost:8443/priority?priority=high"
```

Sample response,

```plaintext
High priority stream # takes 1 second to respond
Stream priority demonstration. Check with an HTTP/2 client.
```

### 3.4. Server Push

- **What to expect**: The server proactively pushes associated resources to the client.
- **How to test**:
  - Using `curl`, `curl -k --http2 -v https://localhost:8443/push`
  - Check the verbose output to see style.css and script.js being pushed.
  - Alternatively, test in a browser (e.g., Chrome) by visiting <https://localhost:8443/push> and viewing the Network tab in Developer Tools.

## Advanced Testing with `h2load`

- What is `h2load`?: A benchmarking tool for HTTP/2.
- How to test,
  - Install `h2load`, for example, with `sudo apt install nghttp2`
  - Run a test to observe multiplexing and server performance, `h2load -n 100 -c 10 https://localhost:8443/multiplex`
  - Observe the results, including how requests are handled concurrently.

## Debugging and Observing HTTP/2 Behavior

- Tools for deeper analysis,
  - Wireshark:
    - Monitor network packets and observe header compression (HPACK).
  - Browser DevTools:
    - View HTTP/2 requests and responses under the Network tab.
  - Log Outputs:
    - Use log.Printf in the server code to trace connections and request handling.

## Additional Notes

Use `-k` with `curl` to bypass certificate warnings during local testing.

Replace `localhost` with your server's hostname or IP when deploying in a production environment.
