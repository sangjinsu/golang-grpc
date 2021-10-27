
# gRPC Golang Master Class

section 06 gRPC Server Streaming

## What's a Server Streaming API?

- Server Steaming RPC API are a New kind API enabled thanks to HTTP/2
- The client will send one message to the server 
- The client will receive many responses from the server, possibly an infinite number

- Streaming Server are well suited for when the server needs to send a lot of data
- or When the server needs to PUSH data to the client without having the client request for more 



- In gRPC Server streaming Calls are defined using the keyword stream

  
