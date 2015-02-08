# Buzzer

Buzzer takes advantage of resolving one hostname into multiple IPs.
`net.Dial` and `net.DialTimeout` in Go only try single IP address
and fail if it's unavailable. Buzzer tries them all to ensure
high availability and resiliency of your application.

## Usage

`buzzer.Dial` and `buzzer.DialTimeout` are direct replacements
for `net.Dial` and `net.DialTimeout`.

Here's an example how to use it for http requests:

```go
package main

import (
    "log"
    "io"
    "os"
    "http"
    "github.com/bobrik/buzzer"
)

func main() {
    c := &http.Client{
        Transport: &http.Transport{
            Dial: buzzer.Dial,
        },
    }

    r, err := c.Get("http://google.com/")
    if err != nil {
        log.Fatal(err)
    }

    defer r.Body.Close()

    io.Copy(os.Stdout, r.Body)
}
```

Note that `buzzer` treats timeout for `buzzer.DialTimeout` as
timeout for single connection attempt. If your domain resolves
to 10 different IPs and all of them are unavailable, call
would block for 10x requested time.

## License

The MIT License (MIT)
