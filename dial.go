package buzzer

import (
	"net"
	"syscall"
	"time"
)

// used in dial_test.go
var resolver = net.LookupIP

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

			switch e := err.(type) {
			case *net.OpError:
				if e.Err == syscall.ECONNREFUSED || e.Temporary() {
					continue
				}
			default:
				return conn, err
			}
		}

		return conn, err
	}

	return nil, lastErr
}

func Dial(network, address string) (net.Conn, error) {
	return dial(network, address, net.Dialer{})
}

func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return dial(network, address, net.Dialer{Timeout: timeout})
}
