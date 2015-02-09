// Buzzer takes advantage of resolving one hostname into multiple IPs.
// net.Dial and net.DialTimeout in Go only try single IP address
// and fail if it's unavailable. Buzzer tries them all to ensure
// high availability and resiliency of your application.
package buzzer

import (
	"net"
	"time"
)

// used in dial_test.go
var resolver = net.LookupIP

// dial tries to connect to the first available IP address
// if network is "tcp" and hostname resolves to several IPs
func dial(network, address string, dialer net.Dialer) (net.Conn, error) {
	if network != "tcp" {
		return net.Dial(network, address)
	}

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	ips, err := resolver(host)
	if err != nil {
		return nil, err
	}

	var lastErr error
	for _, ip := range ips {
		conn, err := dialer.Dial(network, net.JoinHostPort(ip.String(), port))
		if err != nil {
			lastErr = err
			continue
		}

		return conn, err
	}

	return nil, lastErr
}

// Dial works just like net.Dial, but also tries to connect
// to any available IP if several are available
func Dial(network, address string) (net.Conn, error) {
	return dial(network, address, net.Dialer{})
}

// DialTimeout works just like net.DialTimeout, but also tries
// to connect to any available IP if several are available
func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return dial(network, address, net.Dialer{Timeout: timeout})
}
