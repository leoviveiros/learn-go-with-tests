package racer

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRacer(t *testing.T) {
	t.Run("returns the faster url", func(t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0)
	
		defer slowServer.Close()
		defer fastServer.Close()
	
		slowUrl := slowServer.URL
		fastUrl := fastServer.URL
	
		want := fastUrl
		got, err := Racer(slowUrl, fastUrl)

		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
	
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns an error after 10s", func(t *testing.T) {
		server := makeDelayedServer(time.Millisecond * 25)

		defer server.Close()

		_, err := ConfigurableRacer(server.URL, server.URL, time.Millisecond * 20)

		if err == nil {
			t.Error("expected an error")
		}
	})
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(delay)
        w.WriteHeader(http.StatusOK)
    }))
}
