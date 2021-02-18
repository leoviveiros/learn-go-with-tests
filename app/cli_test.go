package poker

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestCLI(t *testing.T) {
    // var dummyBlindAlerter = &SpyBlindAlerter{}
    // var dummyPlayerStore = &StubPlayerStore{}
    var dummyStdOut = &bytes.Buffer{}
    
    t.Run("record chris win from user input", func(t *testing.T) {
        in := strings.NewReader("5\nChris wins\n")
        game := &GameSpy{}
        
        cli := NewCLI(in, dummyStdOut, game)
        cli.PlayPoker()

        if game.FinishedWith != "Chris" {
            t.Errorf("wanted Finish with Chris but got %v", game.FinishedWith)
        }
    })

    t.Run("record cleo win from user input", func(t *testing.T) {
        in := strings.NewReader("5\nCleo wins\n")
        game := &GameSpy{}

        cli := NewCLI(in, dummyStdOut, game)
        cli.PlayPoker()

        if game.FinishedWith != "Cleo" {
            t.Errorf("wanted Finish with Cleo but got %v", game.FinishedWith)
        }
    })

    t.Run("it schedules printing of blind values", func(t *testing.T) {
        in := strings.NewReader("5\nChris wins\n")
        playerStore := &StubPlayerStore{}
        blindAlerter := &SpyBlindAlerter{}

        game := NewGame(blindAlerter, playerStore)

        cli := NewCLI(in, dummyStdOut, game)
        cli.PlayPoker()

        cases := []ScheduledAlert{
            {0 * time.Second, 100},
            {10 * time.Minute, 200},
            {20 * time.Minute, 300},
            {30 * time.Minute, 400},
            {40 * time.Minute, 500},
            {50 * time.Minute, 600},
            {60 * time.Minute, 800},
            {70 * time.Minute, 1000},
            {80 * time.Minute, 2000},
            {90 * time.Minute, 4000},
            {100 * time.Minute, 8000},
        }

        CheckSchedulingCases(t, cases, blindAlerter)
    })

    t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
        stdout := &bytes.Buffer{}
        in := strings.NewReader("7\n")
        game := &GameSpy{}

        cli := NewCLI(in, stdout, game)
        cli.PlayPoker()

        wantPrompt := PlayerPrompt

        assertMessagesSentToUser(t, stdout, wantPrompt)

        if game.StartedWith != 7 {
            t.Errorf("wanted Start called with 7 but got %d", game.StartedWith)
        }
    })

    t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
        stdout := &bytes.Buffer{}
        in := strings.NewReader("Pies\n")
        game := &GameSpy{}

        cli := NewCLI(in, stdout, game)
        cli.PlayPoker()

        wantPrompt := PlayerPrompt + BadPlayerInputErrMsg

        assertMessagesSentToUser(t, stdout, wantPrompt)

        if game.StartCalled {
            t.Errorf("game should not have started")
        }
    })
}

func assertMessagesSentToUser(t testing.TB, stdout *bytes.Buffer, messages ...string) {
    t.Helper()
    want := strings.Join(messages, "")
    got := stdout.String()
    if got != want {
        t.Errorf("got %q sent to stdout but expected %+v", got, messages)
    }
}


