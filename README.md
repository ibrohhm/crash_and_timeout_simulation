# Crash and Timeout Simulation

Image you have apps that required called partner to served your data, the partner sometimes got unexpected behavior that we cannot control,
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

## Simulation
In this section we will do three simulation
1. partner with no delay
2. partner with random delay but no timeout set in the server
3. partner with random delay with timeout set in the server

### Case 1
Every case always have happy case and this is it, our partner service have good spec and never got delay everytime we request.
to make this possible you need to change the delay on partner code from `delay := time.Duration(rand.Intn(11))` to `delay := time.Duration(0)` ([ref](https://github.com/ibrohhm/crash_and_timeout_simulation/blob/master/partner/partner.go#L14))

```
client --> server --> partner
```

this is result

https://github.com/user-attachments/assets/14db531c-6ca5-468b-bde9-30fbc84ab1ac

### Case 2
Our partner service have random delay (`delay := time.Duration(rand.Intn(11))`) and our server service not set the timeout when request to partner service

```
client --> server --> partner (random delay)
```

this is the result

https://github.com/user-attachments/assets/4655345a-e8c9-42b0-ace9-61c2a7511eca

![image](https://github.com/user-attachments/assets/bdb8fae1-bf18-421d-af55-cc16a6bad7e8)

our server got killed in 24 second since the memory usage exceed the MemoryLimit.
This is because the client service continuosly spawn new request using goroutine to call server service, since there's no limit on the number of goroutine being spawend,
the server service request the partner service with hugh number. Because of the partner delay, most of the request running at the same time and consume all available memory then leading to a crash

what happen if we set the timeout request on the server service?

### Case 3
Our partner have random delay but our server set the timeout request.
we need to change the `Timeout` variable in server to some number, let say 3 second `const Timeout = 3` ([ref](https://github.com/ibrohhm/crash_and_timeout_simulation/blob/master/server/server.go#L16))

```
client --> server --[with timeout]--> partner (random delay)
```

this is the result

https://github.com/user-attachments/assets/589afca9-6010-46c8-b68b-f3053f5676c3

If you running the simulation, you'll see our server service not get killed from the exceed memory allocation.
if you look more closely in the logger, the memory_usage of the server is always around 6MB - 12MB (never exceed the 20MB)
this is because the timeout killed the ongoing request and release the memory allocation

<img width="745" alt="image" src="https://github.com/user-attachments/assets/b3e0877f-03ba-4041-b751-7d5cc8df8cfd">

## Summaries
The partner service behaviour is the external thing that we cannot control, we cannot trust the partner to have good behavior.
sometimes it got delay, sometimes we cannot access it doe their internal error or else. the small delay maybe will cause our server down (like the simulation),
so we better to prevent that happen and one way to prevent it is like adding the timeout when request to the server



