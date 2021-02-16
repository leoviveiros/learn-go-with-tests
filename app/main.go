package main

import (
	"log"
	"net/http"
)

func main() {
	server := &PlayerServer{NewInMemoryPlayerStore()}
	err := http.ListenAndServe(":8080", server)

	if err != nil {
		log.Fatalf("could not listen on port 8080 %v", err)
	}
	
}

// curl -X POST http://localhost:8080/players/Pepper
// curl http://localhost:8080/players/Pepper