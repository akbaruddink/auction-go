package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

// Item struct to hold auction item details
type Item struct {
	ID   int
	Name string
	Bids []Bid
	mu   sync.Mutex
}

type ItemRequest struct {
	Name string `json:"name"`
}

type ItemResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// CreateItemHandler handles Item creation
func CreateItemHandler(items *map[int]*Item, itemIDCounter *int, itemMU *sync.Mutex) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var itemR ItemRequest
		if err := json.NewDecoder(r.Body).Decode(&itemR); err != nil || itemR.Name == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		itemMU.Lock()
		defer itemMU.Unlock()

		item := &Item{
			ID:   *itemIDCounter,
			Name: itemR.Name,
			Bids: []Bid{},
		}
		(*items)[*itemIDCounter] = item
		*itemIDCounter++

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ItemResponse{ID: item.ID, Name: item.Name})
	}
}

// GetItemsHandler handles Item(s) retrieval
func GetItemsHandler(items *map[int]*Item, itemMU *sync.Mutex) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		itemMU.Lock()
		defer itemMU.Unlock()

		itemL := make([]ItemResponse, 0, len(*items))
		for _, item := range *items {
			itemL = append(itemL, ItemResponse{ID: item.ID, Name: item.Name})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(itemL)
	}
}

// GetItemWinner handles Item winner retrieval
func GetItemWinner(items *map[int]*Item) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")

		itemID, err := strconv.Atoi(idStr)
		if idStr == "" || err != nil {
			http.Error(w, "Invalid Item ID", http.StatusBadRequest)
			return
		}

		item, found := (*items)[itemID]
		if !found {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		item.mu.Lock()
		defer item.mu.Unlock()

		if len(item.Bids) == 0 {
			http.Error(w, "No winner found", http.StatusNotFound)
			return
		}

		winnerIndex, _ := item.determineWinner()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item.Bids[winnerIndex])
	}
}

// determineWinner determines the winner of the auction Item
func (item *Item) determineWinner() (int, int) {
	highestBid := 0
	winnerIndex := -1

	for {
		for bi := range item.Bids {
			// Fallback to initial bid if current bid is 0
			if item.Bids[bi].CurrentBid == 0 {
				item.Bids[bi].CurrentBid = item.Bids[bi].InitialBid
			}

			for item.Bids[bi].CurrentBid <= highestBid && item.Bids[bi].CurrentBid+item.Bids[bi].BidIncrement <= item.Bids[bi].MaxBid {
				item.Bids[bi].CurrentBid += item.Bids[bi].BidIncrement
			}

			if item.Bids[bi].CurrentBid > highestBid {
				highestBid = item.Bids[bi].CurrentBid
			}
		}

		// Check if highest bid is final
		isFinal := true
		for _, bidder := range item.Bids {
			if bidder.CurrentBid < highestBid && bidder.CurrentBid <= bidder.MaxBid && bidder.CurrentBid+bidder.BidIncrement <= bidder.MaxBid {
				isFinal = false
			}
		}
		if isFinal {
			break
		}
	}

	// Find the winner
	for bi := range item.Bids {
		if item.Bids[bi].CurrentBid == highestBid {
			winnerIndex = bi
			break
		}
	}

	return winnerIndex, highestBid
}
