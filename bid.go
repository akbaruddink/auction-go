package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Bid struct
type Bid struct {
	BidderName   string `json:"bidder_name"`
	InitialBid   int    `json:"initial_bid"`
	MaxBid       int    `json:"max_bid"`
	BidIncrement int    `json:"bid_increment"`
	CurrentBid   int    `json:"current_bid"`
	ItemID       int    `json:"item_id"`
}

type BidRequest struct {
	BidderName   string `json:"bidder_name"`
	InitialBid   int    `json:"initial_bid"`
	MaxBid       int    `json:"max_bid"`
	BidIncrement int    `json:"bid_increment"`
}

// CreateBidHandler handles Bid creation
func CreateBidHandler(items *map[int]*Item) func(http.ResponseWriter, *http.Request) {
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

		var bidR BidRequest
		if err := json.NewDecoder(r.Body).Decode(&bidR); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		item.mu.Lock()
		defer item.mu.Unlock()

		bid := Bid{
			BidderName:   bidR.BidderName,
			InitialBid:   bidR.InitialBid,
			MaxBid:       bidR.MaxBid,
			BidIncrement: bidR.BidIncrement,
			CurrentBid:   bidR.InitialBid,
			ItemID:       itemID,
		}
		item.Bids = append(item.Bids, bid)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bid)
	}
}

// GetBidsHandler handles Bid retrieval
func GetBidsHandler(items *map[int]*Item) func(http.ResponseWriter, *http.Request) {
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item.Bids)
	}
}
