package buzzer

import (
	"net"
	"testing"
	"time"
)

func TestIP(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	_, err = Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	l.Close()

	_, err = Dial("tcp", l.Addr().String())
	if err == nil {
		t.Fatal(err)
	}
}

func TestDisabledIP(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	resolver = func(string) ([]net.IP, error) {
		return []net.IP{
			{192, 0, 2, 0}, // unavailable
			{192, 0, 2, 1}, // unavailable
			{192, 0, 2, 2}, // unavailable
			{127, 0, 0, 1},
		}, nil
	}

	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	addr := net.JoinHostPort("example.com", port)

	_, err = DialTimeout("tcp", addr, 100*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}

	l.Close()

	_, err = DialTimeout("tcp", addr, 100*time.Millisecond)
	if err == nil {
		t.Fatal(err)
	}

	resolver = func(string) ([]net.IP, error) {
		return []net.IP{
			{127, 0, 0, 1},
		}, nil
	}

	_, err = DialTimeout("tcp", addr, 100*time.Millisecond)
	if err == nil {
		t.Fatal(err)
	}

}
