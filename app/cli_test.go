package poker

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

func TestCLI(t *testing.T) {
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

        game := NewTexasHoldem(blindAlerter, playerStore)

        cli := NewCLI(in, dummyStdOut, game)
        cli.PlayPoker()

        dest := os.Stdout

        cases := []ScheduledAlert{
            {0 * time.Second, 100, dest},
            {10 * time.Minute, 200, dest},
            {20 * time.Minute, 300, dest},
            {30 * time.Minute, 400, dest},
            {40 * time.Minute, 500, dest},
            {50 * time.Minute, 600, dest},
            {60 * time.Minute, 800, dest},
            {70 * time.Minute, 1000, dest},
            {80 * time.Minute, 2000, dest},
            {90 * time.Minute, 4000, dest},
            {100 * time.Minute, 8000, dest},
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


