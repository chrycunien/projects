package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect to localhost: %v", err), http.StatusServiceUnavailable)
		return
	}

	// first query the latest block
	peerCount, err := client.PeerCount(context.Background())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read peer count: %v", err), http.StatusServiceUnavailable)
		return
	}

	// define the threshold for minimum number of peers
	var minPeers uint64 = 10

	// compare the peer count with the minimum threshold
	if peerCount < minPeers {
		http.Error(w, fmt.Sprintf("Number of peers is too low: %d", peerCount), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/readiness", readinessHandler)
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
