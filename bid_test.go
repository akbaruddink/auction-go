package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateBidHandler(t *testing.T) {
	items := make(map[int]*Item)

	// Create a new item
	item := &Item{
		ID:   1,
		Name: "Test Item",
		Bids: []Bid{},
	}
	items[item.ID] = item

	// Create a new request with a JSON payload
	bidReq := BidRequest{
		BidderName:   "Test Bidder",
		InitialBid:   100,
		MaxBid:       200,
		BidIncrement: 10,
	}
	reqBody, _ := json.Marshal(bidReq)
	req, _ := http.NewRequest("POST", "/items/{id}/bids", bytes.NewBuffer(reqBody))
	req.SetPathValue("id", "1")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	CreateBidHandler(&items)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	var bid Bid
	json.NewDecoder(rr.Body).Decode(&bid)

	expectedBid := Bid{
		BidderName:   "Test Bidder",
		InitialBid:   100,
		MaxBid:       200,
		BidIncrement: 10,
		CurrentBid:   100,
		ItemID:       1,
	}
	if !reflect.DeepEqual(bid, expectedBid) {
		t.Errorf("Expected bid %+v, but got %+v", expectedBid, bid)
	}
}

func TestCreateBidInvalidID(t *testing.T) {
	items := make(map[int]*Item)

	// Create a new request with an invalid ID
	req, _ := http.NewRequest("GET", "/items/invalid/winner", nil)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	CreateBidHandler(&items)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestCreateBidItemNotFound(t *testing.T) {
	items := make(map[int]*Item)

	// Create a new request with a valid ID but no item
	req, _ := http.NewRequest("GET", "/items/{id}/winner", nil)
	req.SetPathValue("id", "1")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	CreateBidHandler(&items)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, rr.Code)
	}
}

func TestCreateBidInvalidBody(t *testing.T) {
	items := make(map[int]*Item)

	// Create a new item
	item := &Item{
		ID:   1,
		Name: "Test Item",
		Bids: []Bid{},
	}
	items[item.ID] = item

	// Create a new request with an invalid body
	req, _ := http.NewRequest("POST", "/items/{id}/bids", bytes.NewBuffer([]byte("invalid")))
	req.SetPathValue("id", "1")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	CreateBidHandler(&items)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGetBidsHandler(t *testing.T) {
	items := make(map[int]*Item)

	// Create a new item
	item := &Item{
		ID:   1,
		Name: "Test Item",
		Bids: []Bid{
			{
				BidderName:   "Test Bidder 1",
				InitialBid:   100,
				MaxBid:       200,
				BidIncrement: 10,
				CurrentBid:   100,
				ItemID:       1,
			},
			{
				BidderName:   "Test Bidder 2",
				InitialBid:   150,
				MaxBid:       250,
				BidIncrement: 20,
				CurrentBid:   150,
				ItemID:       1,
			},
		},
	}
	items[item.ID] = item

	// Create a new request
	req, _ := http.NewRequest("GET", "/items/{id}/bids", nil)
	req.SetPathValue("id", "1")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	GetBidsHandler(&items)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	var bids []Bid
	json.NewDecoder(rr.Body).Decode(&bids)

	expectedBids := []Bid{
		{
			BidderName:   "Test Bidder 1",
			InitialBid:   100,
			MaxBid:       200,
			BidIncrement: 10,
			CurrentBid:   100,
			ItemID:       1,
		},
		{
			BidderName:   "Test Bidder 2",
			InitialBid:   150,
			MaxBid:       250,
			BidIncrement: 20,
			CurrentBid:   150,
			ItemID:       1,
		},
	}
	if !reflect.DeepEqual(bids, expectedBids) {
		t.Errorf("Expected bids %+v, but got %+v", expectedBids, bids)
	}
}

func TestGetBidsInvalidID(t *testing.T) {
	items := make(map[int]*Item)

	// Create a new request with an invalid ID
	req, _ := http.NewRequest("GET", "/items/invalid/winner", nil)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	GetBidsHandler(&items)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGetBidsItemNotFound(t *testing.T) {
	items := make(map[int]*Item)

	// Create a new request with a valid ID but no item
	req, _ := http.NewRequest("GET", "/items/{id}/winner", nil)
	req.SetPathValue("id", "1")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	GetBidsHandler(&items)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, rr.Code)
	}
}
