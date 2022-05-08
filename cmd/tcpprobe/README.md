## purpose

tiny go application to explore building different kinds of network probes atop the go standard library.

## status

*	supports TCP
*	supports TLS
*	failed probes are communicated in terrible way (panic on probe fail!)

## usage

### build it

```
go build -o main main.go
```

### run a TCP probe

```
./main -url=tcp://example.com:80
```

### run a TLS probe, using root CAs trusted by the operating system

```
./main -url=tls://example.com:443
```

### run a TLS probe, distrusting all root CAs

```
./main -url=tls://example.com:443 -tls-trust-no-one
```

### run it while tracing linux kernel `tcp_sendmsg`:

```
sudo ./tcpperf.sh
```
