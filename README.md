# Crash and Timeout Simulation

Image you have apps that required called partner to server your data, the partner sometimes got unexpected behavior that we cannot control,
let say it's random delay everytime you request from the partner.
This repo is focus on the simulation for your server to handle this partner behavior

To simulate this we will create three service (client, server, partner) using golang

```
client --> server --> partner
```

- client: call the server with go routine
- server: the server will forward the request from client to partner, act as middleware
- partner: simple hello world golang with random timeout

## Partner Service
Partner service is simple http call with endpoint `get /data` that served in port 8081 with random delay

The partner service only have `get /data` endpoint with response _Hello from Partner Service_
with generate random delay everytime request the data (1-10 second delay).
The partner also have logging to show the _delay_set_ and _time_ when request occur. So we can monitor the request well

See the implementation: ([partner service](https://github.com/ibrohhm/crash_and_timeout_simulation/blob/master/partner/partner.go))

How to run: `go run partner.go`
