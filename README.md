# HTTP Showcase

This repository contains small, self-contained server implementations that demonstrate key features of HTTP/1.1, HTTP/2, and HTTP/3. The goal is to explore the evolution of the HTTP protocol and showcase how each version improves on its predecessors while maintaining backward compatibility.  

## Features  

### HTTP/1.1  

- Persistent connections  
- Chunked transfer encoding  
- Host header handling  
- Pipelining (optional)  

### HTTP/2  

- Multiplexed streams  
- Header compression (HPACK)  
- Stream prioritization  
- Server push  

### HTTP/3  

- QUIC-based transport (UDP)  
- Improved latency and connection establishment  
- Multiplexing without head-of-line blocking  

## Requirements  

- Go (1.18 or later recommended)  
- A modern browser or tools like `curl`/`Postman` to test server behavior  

## Getting Started  

1. Clone the repository:  

   ```bash  
   git clone https://github.com/yourusername/http-showcase.git  
   cd http-showcase  
   ```  

2. Run the server for a specific HTTP version:  
   - HTTP/1.1:  

     ```bash  
     go run ./http1.1/server.go  
     ```  

   - HTTP/2:  

     ```bash  
     go run ./http2/server.go  
     ```  

   - HTTP/3:  

     ```bash  
     go run ./http3/server.go  
     ```  

3. Test the server using a browser or tool.  

## Project Structure  

```text
http-showcase-servers/  
├── http1.1/  
│   └── server.go  
├── http2/  
│   └── server.go  
└── http3/  
    └── server.go  
```  
