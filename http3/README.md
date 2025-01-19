# # HTTP/3 Showcase

### Multiplexed Streams

Run the `main.go` file to start the HTTP/3 server.

```sh
go run main.go
```

Follow the below steps to use `Toxiproxy` to test the HTTP/3 server.

### Toxiproxy

```sh
docker run -it -p 8474:8474 -p 1080:1080 ghcr.io/shopify/toxiproxy
```

You should get these logs.

```shell
{"level":"info","version":"2.11.0","caller":"server.go:78","time":"2025-01-19T07:59:26Z","message":"Starting Toxiproxy"}
{"level":"info","address":"0.0.0.0:8474","caller":"api.go:57","time":"2025-01-19T07:59:26Z","message":"Starting Toxiproxy HTTP server"}
{"level":"info","name":"http3_proxy","listen":"127.0.0.1:1080","upstream":"127.0.0.1:8443","caller":"proxy.go:118","time":"2025-01-19T07:59:40Z","message":"Started proxy"}
```

Install `toxiproxy-cli` on your local machine.

For `MacOs`, you can run the following command to install `toxiproxy-cli`.

```sh
brew install toxiproxy
```

Create a new proxy.

```shell
toxiproxy-cli create --listen 127.0.0.1:1080 --upstream 127.0.0.1:8443 http3_proxy
# Created new proxy http3_proxy
```

Add a packet loss rule to the proxy.

```shell
toxiproxy-cli toxic add --type=limit_data --toxicName=packet_loss --toxicity=0.2 http3_proxy
# Added downstream limit_data toxic 'packet_loss' on proxy 'http3_proxy'
```

### Client

Start the client to make a request to the HTTP/3 server.

```sh
go run client/client.go
```

You should see the following logs for the server.

```plaintext
2025/01/19 13:07:08 HTTP/3 Server is running on https://localhost:8443
2025/01/19 13:07:12 Received request: GET / from 127.0.0.1:61612
2025/01/19 13:07:12 Received request: GET /static/style.css from 127.0.0.1:61612
2025/01/19 13:07:12 Simulating packet loss
2025/01/19 13:07:12 Received request: GET /push from 127.0.0.1:61612
2025/01/19 13:07:12 Hello from HTTP/3 with Server Push!
2025/01/19 13:07:12 Received request: GET /stream6 from 127.0.0.1:61612
2025/01/19 13:07:12 Received request: GET /stream5 from 127.0.0.1:61612
2025/01/19 13:07:12 Received request: GET /stream8 from 127.0.0.1:61612
2025/01/19 13:07:12 Received request: GET /stream3 from 127.0.0.1:61612
2025/01/19 13:07:12 Received request: GET /stream10 from 127.0.0.1:61612
2025/01/19 13:07:12 Received request: GET /stream1 from 127.0.0.1:61612
2025/01/19 13:07:12 Received request: GET /stream7 from 127.0.0.1:61612
2025/01/19 13:07:12 Received request: GET /stream9 from 127.0.0.1:61612
2025/01/19 13:07:12 Simulating packet loss
2025/01/19 13:07:12 Received request: GET /stream2 from 127.0.0.1:61612
2025/01/19 13:07:23 Stream 5 responded after 11 ms
2025/01/19 13:07:25 Stream 7 responded after 13 ms
2025/01/19 13:07:25 Stream 2 responded after 13 ms
2025/01/19 13:07:27 Stream 3 responded after 15 ms
2025/01/19 13:07:33 Stream 9 responded after 21 ms
2025/01/19 13:07:35 Stream 8 responded after 23 ms
2025/01/19 13:07:35 Stream 1 responded after 23 ms
2025/01/19 13:07:35 Stream 6 responded after 23 ms
2025/01/19 13:07:41 Stream 10 responded after 29 ms
```

Sample client logs.

```plaintext
Response from https://localhost:8443:
Status: 200 OK
Body:
<html>
<head>
    <link rel='stylesheet' href='/style.css'>
</head>
<body><h1>Welcome to HTTP/3</h1>
<script src='/script.js'></script>
</body>
</html>

Response from https://localhost:8443/static/style.css:
Status: 200 OK
Body:
body {
    font-family: Arial, sans-serif;
}

Response from https://localhost:8443/static/script.js:
Status: 200 OK
Body:

Response from https://localhost:8443/push:
Status: 200 OK
Body:

Response from https://localhost:8443/stream4?delay=17: 
Response from https://localhost:8443/stream5?delay=11: 
Response from https://localhost:8443/stream7?delay=13: 
Response from https://localhost:8443/stream2?delay=13: 
Response from https://localhost:8443/stream3?delay=15: 
Response from https://localhost:8443/stream9?delay=21: 
Response from https://localhost:8443/stream8?delay=23: 
Response from https://localhost:8443/stream1?delay=23: 
Response from https://localhost:8443/stream6?delay=23: 
Response from https://localhost:8443/stream10?delay=29:
```
