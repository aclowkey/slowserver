# slowserver
A simple HTTP server which will be as slow as you want!

## Motivation
The reason the project was create is to test how our proxy handle blue/green deployment switch.

## Usage
The server answers to any HTTP request at `<host>/timeout`.
You can optionally provide a `?timeout=5m` parameter in the request.
There is a safety timeout flag, so the server will wait the maximum of the param, and the flag 

```
Usage of slowserver:
  -listen string
        address and port to listen (default ":4211")
  -max-timeout string
        maximum timeout (default "10m")
```

## Example
```sh
$ docker run -p 8080:8080 aclowkey/slowserver -listen :8080 -max-timeout 1m
$ time curl localhost:8080/timeout?duration=15s
real    0m15.027s
user    0m0.019s
sys     0m0.005s
```
