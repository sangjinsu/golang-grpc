
# gRPC Golang Master Class

section 07 gRPC Client Streaming

## What's a Client Streaming API?

- Client Streaming RPC API are a NEW kind API enabled thanks to HTTP/2

- The client will send many message to the server and will receive one response from the server 

- Streaming Client are well suited for When the client needs to send a lot of data

  - When the Server processing is expensive and should happen as the client sends data

  - When the client needs to PUSH data to the server without really expecting a response

    

- In gRPC Client Streaming Calls are defined using the keyword "stream"
