package racer

import (
	"fmt"
	"net/http"
	"time"
)

var tenSecondsTimeout = time.Second * 10

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, err error) {
	select {
		case <-ping(a):
			return a, nil
		case <-ping(b):
			return b, nil
		case <-time.After(timeout):
			return "", fmt.Errorf("timed out waiting for %q and %q", a, b)
	}
}

func Racer(a, b string) (winner string, err error) {
	return ConfigurableRacer(a, b, tenSecondsTimeout)
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})

	go getUrl(url, ch)

	return ch
}

func getUrl(url string, ch chan struct{}) {
	http.Get(url)
	close(ch)
}