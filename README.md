# tcp-proxy
A simple tcp connection proxy.

## Install

```shell
go install github.com/rselph/tcpproxy@latest
```

## Run

```shell
> tcpproxy -h
Usage of tcpproxy:
  -forward string
        Address to forward requests to
  -listen string
        Address to listen on
  -protocol string
        Protocol to use (tcp, tcp4, tcp6) (default "tcp4")
```
