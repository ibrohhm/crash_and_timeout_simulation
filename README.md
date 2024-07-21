# Crash and Timeout Simulation

Image you have apps that required called partner to server your data, the partner sometimes got unexpected behavior that we cannot control,
let say it's random delay everytime you request from the partner. It's very tiny detail but if we are not handle the partner request well, it will causing our server down.
This repo is focus on the simulation for your server to handle this partner behavior

To simulate this we will create three service (client, server, partner) using golang

```
client --> server --> partner
```

- client: call the server with go routine
- server: the server will forward the request from client to partner, act as middleware
- partner: simple hello world golang with random timeout

## Partner Service
Partner service is simple http call with random delay

The partner service only have `get /data` endpoint with response _Hello from Partner Service_
with generate random delay everytime request the data (1-10 second delay).
The partner also have logging to show the _delay_set_ and _time_ when request occur. So we can monitor the request well

See the implementation: ([partner service](https://github.com/ibrohhm/crash_and_timeout_simulation/blob/master/partner/partner.go))

How to run: `go run partner.go`

## Server Service
Server service is your internal service to handle the request from client. 
To simulate the crash, we need to set the memory limit allocation (`MemoryLimit`) so we can simulate the crash without crashing your laptop.
When running the server, it will checking the memory usage in every 1 second (`getMemoryUsage`) and the memory usage is exceed the `MemoryLimit` we will stop the server.
The service also have logging to show the _method_, _url_, _latency_, _status_, _error_, and _memory_usage_

See the implementation: [server service](https://github.com/ibrohhm/crash_and_timeout_simulation/blob/master/server/server.go)

How to run: `go run server.go`

## Client Service
Client service is simple golang apps that will do 100 request in 1 second to the server with go routine. I choose to create client service instead of using load test application like _JMeter_, so we can see the logger for every request

See the implementation: [client service](https://github.com/ibrohhm/crash_and_timeout_simulation/blob/master/client/client.go)

How to run: `go run client.go`

