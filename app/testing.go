package poker

import (
	"fmt"
	"io"
	"testing"
	"time"
)

var DummyGame = &GameSpy{}

type ScheduledAlert struct {
	At     time.Duration
	Amount int
    To io.Writer
}

func (s ScheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.Amount, s.At)
}

type SpyBlindAlerter struct {
	alerts []ScheduledAlert
}

func (s *SpyBlindAlerter) ScheduleAlertAt(at time.Duration, amount int, to io.Writer) {
	s.alerts = append(s.alerts, ScheduledAlert{at, amount, to})
}

type StubPlayerStore struct {
    scores map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
    score := s.scores[name]
    return score
}

func (s *StubPlayerStore) RecordWin(name string) {
    s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
    return s.league
}

type GameSpy struct {
    StartedWith  int
	StartCalled bool
    BlindAlert  []byte

    FinishedCalled   bool
    FinishedWith string
}

func (g *GameSpy) Start(numberOfPlayers int, out io.Writer) {
    g.StartedWith = numberOfPlayers
	g.StartCalled = true
    out.Write(g.BlindAlert)
}

func (g *GameSpy) Finish(winner string) {
    g.FinishedWith = winner
}

func AssertScheduledAlert(t testing.TB, got, want ScheduledAlert) {
    t.Helper()
    if got.Amount != want.Amount {
        t.Errorf("got amount %d, want %d", got.Amount, want.Amount)
    }

    if got.At != want.At {
        t.Errorf("got scheduled time of %v, want %v", got.At, want.At)
    }
}

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
    t.Helper()

    if len(store.winCalls) != 1 {
        t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
    }

    if store.winCalls[0] != winner {
        t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
    }
}

func CheckSchedulingCases(t *testing.T, cases []ScheduledAlert, blindAlerter *SpyBlindAlerter) {
	t.Helper()

	for i, want := range cases {
		t.Run(fmt.Sprint(want), func(t *testing.T) {

			if len(blindAlerter.alerts) <= i {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}

			got := blindAlerter.alerts[i]
			AssertScheduledAlert(t, got, want)
		})
	}
}