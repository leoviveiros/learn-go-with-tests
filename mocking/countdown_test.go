package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)


func TestCountdown(t *testing.T) {

	t.Run("3 prints to Go", func(t *testing.T) {
		buffer := bytes.Buffer{}
		operationSpy := CountdownOperationsSpy{}

		Countdown(&buffer, &operationSpy)

		got := buffer.String()
		want := "3\n2\n1\nGo!"

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})

	t.Run("sleep before every print", func(t *testing.T) {
		operationSpy := CountdownOperationsSpy{}

		Countdown(&operationSpy, &operationSpy)

		want := []string{
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
		}

		if !reflect.DeepEqual(want, operationSpy.Calls) {
			t.Errorf("wanted calls %v got %v", want, operationSpy.Calls)
		}
	})
}

func TestConfigurableSleeper(t *testing.T) {
    sleepTime := 5 * time.Second

    spyTime := &SpyTime{}
    sleeper := ConfigurableSleeper{sleepTime, spyTime.Sleep}
    sleeper.Sleep()

    if spyTime.durationSlept != sleepTime {
        t.Errorf("should have slept for %v but slept for %v", sleepTime, spyTime.durationSlept)
    }
}

type CountdownOperationsSpy struct {
    Calls []string
}

func (s *CountdownOperationsSpy) Sleep() {
    s.Calls = append(s.Calls, sleep)
}

func (s *CountdownOperationsSpy) Write(p []byte) (n int, err error) {
    s.Calls = append(s.Calls, write)
    return
}

type SpyTime struct {
    durationSlept time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
    s.durationSlept = duration
}

const write = "write"
const sleep = "sleep"