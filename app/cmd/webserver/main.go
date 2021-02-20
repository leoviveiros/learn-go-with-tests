package main

import (
	"log"
	"net/http"

	poker "learn-go-with-tests/app"
)

const dbFileName = "game.db.json"

func main() {
    store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

    if err != nil {
        log.Fatal(err)
    }
    defer close()
    
    game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), store)
    server, err := poker.NewPlayerServer(store, game)

    if err != nil {
        log.Fatal("problem creating player server", err)
    }

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}

// curl -X POST http://localhost:5000/players/Pepper
// curl http://localhost:5000/players/Pepper