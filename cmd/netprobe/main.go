package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"
)

type URLValue struct {
	URL *url.URL
}

func (v URLValue) String() string {
	if v.URL != nil {
		return v.URL.String()
	}
	return ""
}

func (v URLValue) Set(s string) error {
	if u, err := url.Parse(s); err != nil {
		return err
	} else {
		*v.URL = *u
	}
	return nil
}

func main() {
	// TODO: better CLI UX to communicate probe status & failures
	handleErr := func(err error) {
		if err == nil {
			return
		}
		panic(err.Error())
	}

	var u = &url.URL{}
	var uflag = URLValue{URL: u}

	fs := flag.NewFlagSet("tcpprobe", flag.ExitOnError)
	fs.Var(&uflag, "url", "remote url [scheme://]host[:port]")

	var trustNoOne bool
	fs.BoolVar(&trustNoOne, "tls-trust-no-one", false, "TLS: trust no root CAs at all?")

	err := fs.Parse(os.Args[1:])
	handleErr(err)

	// TODO: expose timeout as flag
	timeout := 5 * time.Second

	// If we attempt to connect to a remote host & port that is not accepting TCP
	// connections, Dial will block and then time out.
	var conn net.Conn
	switch u.Scheme {
	case "tcp":
		conn, err = net.DialTimeout(u.Scheme, u.Host, timeout)
	case "tls":
		dialer := &net.Dialer{
			Timeout: timeout,
		}
		// TODO: expose flags to restrict TLS protocol version
		// TODO: expose flags to restrict TLS cipher suites
		tlsCfg := &tls.Config{}
		// TODO: expose flag to read custom root CAs from filesystem path
		if trustNoOne {
			tlsCfg.RootCAs = x509.NewCertPool() // empty!
		}
		conn, err = tls.DialWithDialer(dialer, "tcp", u.Host, tlsCfg)
	}

	handleErr(err)

	fmt.Printf("local addr: %s\n", conn.LocalAddr().String())
	fmt.Printf("remote addr: %s\n", conn.RemoteAddr().String())

	if tlsConn, ok := conn.(*tls.Conn); ok {
		fmt.Printf("got tls conn\n")
		state := tlsConn.ConnectionState()
		fmt.Printf("TLS ConnectionState.HandshakeComplete %v\n", state.HandshakeComplete)
		fmt.Printf("TLS ConnectionState.Version %s\n", TLSVersionName(state.Version))
		fmt.Printf("TLS ConnectionState.CipherSuite %s\n", tls.CipherSuiteName(state.CipherSuite))
	}

	handleErr(err)
}

func TLSVersionName(v uint16) string {
	switch v {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "unknown"
	}
}
