#!/bin/bash

set -e

# ref: https://www.brendangregg.com/perf.html

perf probe --add tcp_sendmsg
perf record -e probe:tcp_sendmsg -ag -- ./main -url=tls://example.com:443
perf probe --del tcp_sendmsg
perf report
