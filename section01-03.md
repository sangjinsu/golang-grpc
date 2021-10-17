# gRPC Golang Master Class

section01 ~ 03

- Microservices are built in different language and encompass a function of business

### Building an API is hard

- Need to think about data model

  json, xml, binary

- Need to think about the endpoint

- Need to think about how to invoke it and handle errors

- Need to think about efficiency of the API

  etc....

### What is an API?

- send Request and fetch Response

### What is gRPC?

- a free and open source framework 
- part of the cloud native computation foundation like Docker and K8S
- define Request and Response for RPC(Remote Procedure Calls)

- <u>it is **modern, fast and efficient**, build on top of **HTTP/2, low  latency**, supports **streaming, language independent**, and makes it super easy to plug in **authentication, load balancing, logging and monitoring**</u>

- An RPC is a Remote Procedure Call

  ![gRPC](README.assets/grpc.svg)

### Why Protocol Buffers?

- Code can be generated for pretty much any language
- Data is binary and efficiently serialized (small payloads)
- Very convenient for transporting a lot of data 
- Protocol Buffers allows for easy API evolution using rules

### Why Should I Learn It?

- Many companies have embraced it fully in Production
- gRPC is the future of microservices API and mobile-server API

# gRPC

Protocol Buffers is used to define the 

- Message (data, req and res)
- Service (Service name and RPC endpoints)

### Efficiency of Protocol Buffers over JSON

- gRPC uses Protocol Buffers for communications
- Protocol Buffer use less memory than json
- Parsing Json is actually CPU intensive
- Parsing Protocol Buffers is less CPU intensive 

## HTTP/2

- gRPC leverages HTTP/2 as a backbone for communications
- https://imagekit.io/demo/http2-vs-http1

- HTTP/2 is the newer standard for internet communications

### HTTP/1.1

- released in 1997
- new TCP connection to a server at each request
- it does not compress headers which are plaintext
- req / res mechanism 

### HTTP/2

- released in 2015
- supports multiplexing
  - The client and server can push messages in parallel over the same TCP connection
  - this greatly reduces latency

- HTTP 2 supports server push
  - Servers can push streams (multiple messages) for one request from the client 
  - this saves round trips (latency)

- supports header compression
  - text based headers can be compressed
  - much less impact on the packet size

- binary
  * Protocol buffers is a binary protocol and makes it a great match for HTTP2

- HTTP/2 is secure 
  - SSL  is not required but recommend by default

## 4 Types of API in gRPC

1. Unary
   - a traditional API looks like HTTP REST
2. Server Streaming
3. Client Streaming
4. Bi Directional Streaming

## Scalability in gRPC

- gRPC Servers are asynchronous by default

- This means they do not block threads on request

- Therefore each gRPC server can serve millions of requests in parallel

  

- gRPC Clients can be asynchronous or synchronous (blocking)
- The client decides which model works best for the performance needs
- gRPC Clients can perform client side load balancing 

## Security in gRPC (SSL)

- gRPC strongly advocates for you to use SSL (encryption over the wire) in your API
- This means that gRPC has security as a first class citizen
- Each language will provide an API to load gRPC with the required certificates and provide encryption capability out of the box
- Additionally using Interceptors, we can also provide authentication

## gRPC vs REST

| gRPC                                                     | REST                                                         |
| -------------------------------------------------------- | ------------------------------------------------------------ |
| Protocol Buffers - smaller, faster                       | JSON - text based, slower, bigger                            |
| HTTP/2 lower latency                                     | HTTP1.1 higher latency                                       |
| Bidirectional and Async                                  | Client => Server requests only                               |
| Stream Support                                           | Request / Response support only                              |
| API Oriented - "What", no constraints                    | CRUD Oriented - GET POST PUT DELETE                          |
| Code generation through Protocol Buffers in any language | Code generation through OpenAPI / Swagger                    |
| RPC Based - gRPC does the plumbing(파이프?) for us       | HTTP verbs based - we have to write the plumbing or use a 3rd party library |

## Summary

- Easy code definition in over 11 languages
- Use a modern, low latency HTTP/2 transport mechanism
- SSL Security is built in 
- Support for streaming APIs for maximum performance
- gRPC is API oriented, instead of Resource Oriented like REST 
