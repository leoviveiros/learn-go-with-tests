package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)


func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
        map[string]int{
            "Pepper": 20,
            "Floyd":  10,
        },
		nil,
		nil,
    }
	server, _ := NewPlayerServer(&store, DummyGame)

    t.Run("returns Pepper's score", func(t *testing.T) {
        request := newGetScoreRequest("Pepper")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
        assertResponseBody(t, response.Body.String(), "20")
    })

    t.Run("returns Floyd's score", func(t *testing.T) {
        request := newGetScoreRequest("Floyd")
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
        assertResponseBody(t, response.Body.String(), "10")
    })

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
    store := StubPlayerStore{
        map[string]int{},
		nil,
		nil,
    }
    server, _ := NewPlayerServer(&store, DummyGame)

	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Pepper"
	
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()
	
		server.ServeHTTP(response, request)
	
		assertStatus(t, response, http.StatusAccepted)
        AssertPlayerWin(t, &store, player)
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as JSON", func(t *testing.T) {
        wantedLeague := []Player{
            {"Cleo", 32},
            {"Chris", 20},
            {"Tiest", 14},
        }

        store := StubPlayerStore{nil, nil, wantedLeague}
        server, _ := NewPlayerServer(&store, DummyGame)

        request := newLeagueRequest()
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        got := getLeagueFromResponse(t, response.Body)
        assertStatus(t, response, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, jsonContentType)
    })
}

func TestGame(t *testing.T) {
    var dummyPlayerStore = &StubPlayerStore{}
    
    t.Run("GET /game returns 200", func(t *testing.T) {
        server := mustMakePlayerServer(t, &StubPlayerStore{}, DummyGame)

        request := newGameRequest()
        response := httptest.NewRecorder()

        server.ServeHTTP(response, request)

        assertStatus(t, response, http.StatusOK)
    })

    t.Run("start a game with 3 players, send some blind alerts down WS and declare Ruth the winner", func(t *testing.T) {
        wantedBlindAlert := "Blind is 100"
        winner := "Ruth"
    
        game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
        server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
        ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")
    
        defer server.Close()
        defer ws.Close()
    
        writeWSMessage(t, ws, "3")
        writeWSMessage(t, ws, winner)
    
        time.Sleep(10 * time.Millisecond)
        assertGameStartedWith(t, game, 3)
        assertFinishCalledWith(t, game, winner)
        within(t, 10 * time.Millisecond, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })
    })
}

func newGetScoreRequest(name string) *http.Request {
    req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
    return req
}

func newPostWinRequest(name string) *http.Request {
    req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
    return req
}

func newGameRequest() *http.Request {
    req, _ := http.NewRequest(http.MethodGet, "/game", nil)
    return req
}

func assertResponseBody(t testing.TB, got, want string) {
    t.Helper()
    if got != want {
        t.Errorf("response body is wrong, got %q want %q", got, want)
    }
}

func assertStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
    t.Helper()
    if response.Code != want {
        t.Errorf("did not get correct status, got %d, want %d", response.Code, want)
    }
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
    t.Helper()
    err := json.NewDecoder(body).Decode(&league)

    if err != nil {
        t.Fatalf("Unable to parse response from server %q into slice of Player, '%v'", body, err)
    }

    return
}

func assertLeague(t testing.TB, got, want []Player) {
    t.Helper()
    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v want %v", got, want)
    }
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
    t.Helper()
    if response.Result().Header.Get("content-type") != want {
        t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
    }
}

func assertGameStartedWith(t testing.TB, game *GameSpy, want int) {
    t.Helper()
    if game.StartedWith != want {
        t.Errorf("game did not start with %d, got %d", want, game.StartedWith)
    }

}

func assertFinishCalledWith(t testing.TB, game *GameSpy, winner string) {
    t.Helper()

    passed := retryUntil(500*time.Millisecond, func() bool {
        return game.FinishedWith == winner
    })

    if !passed {
        t.Errorf("expected finish called with %q but got %q", winner, game.FinishedWith)
    }
}

func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
    _, msg, _ := ws.ReadMessage()
    if string(msg) != want {
        t.Errorf(`got "%s", want "%s"`, string(msg), want)
    }
}

func newLeagueRequest() *http.Request {
    req, _ := http.NewRequest(http.MethodGet, "/league", nil)
    return req
}

func mustMakePlayerServer(t *testing.T, store PlayerStore, game Game) *PlayerServer {
    server, err := NewPlayerServer(store, game)
    if err != nil {
        t.Fatal("problem creating player server", err)
    }
    return server
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
    ws, _, err := websocket.DefaultDialer.Dial(url, nil)

    if err != nil {
        t.Fatalf("could not open a ws connection on %s %v", url, err)
    }

    return ws
}

func writeWSMessage(t testing.TB, conn *websocket.Conn, message string) {
    t.Helper()
    if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
        t.Fatalf("could not send message over ws connection %v", err)
    }
}

func within(t testing.TB, d time.Duration, assert func()) {
    t.Helper()

    done := make(chan struct{}, 1)

    go func() {
        assert()
        done <- struct{}{}
    }()

    select {
    case <-time.After(d):
        t.Error("timed out")
    case <-done:
    }
}

func retryUntil(d time.Duration, f func() bool) bool {
    deadline := time.Now().Add(d)
    for time.Now().Before(deadline) {
        if f() {
            return true
        }
    }
    return false
}