package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"sync"
	"testing"
)

func TestCreateItemHandler(t *testing.T) {
	items := make(map[int]*Item)
	itemIDCounter := 1
	var itemMU sync.Mutex

	// Create a new request with a JSON payload
	itemReq := ItemRequest{Name: "Test Item"}
	reqBody, _ := json.Marshal(itemReq)
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(reqBody))

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	CreateItemHandler(&items, &itemIDCounter, &itemMU)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	var itemRes ItemResponse
	json.NewDecoder(rr.Body).Decode(&itemRes)

	expectedItemRes := ItemResponse{ID: 1, Name: "Test Item"}
	if itemRes != expectedItemRes {
		t.Errorf("Expected item response %+v, but got %+v", expectedItemRes, itemRes)
	}
}

func TestCreateItemHandlerInvalidRequest(t *testing.T) {
	items := make(map[int]*Item)
	itemIDCounter := 1
	var itemMU sync.Mutex

	// Create a new request with an invalid JSON payload
	reqBody := []byte(`{"invalid": "payload"}`)
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(reqBody))

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	CreateItemHandler(&items, &itemIDCounter, &itemMU)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGetItemsHandler(t *testing.T) {
	expectedItemRes := []ItemResponse{
		{ID: 1, Name: "Test Item - 1"},
		{ID: 2, Name: "Test Item - 2"},
		{ID: 3, Name: "Test Item - 3"},
	}

	items := make(map[int]*Item)
	var itemMU sync.Mutex

	for _, expectedItem := range expectedItemRes {
		item := Item{
			ID:   expectedItem.ID,
			Name: expectedItem.Name,
			Bids: []Bid{},
		}

		items[expectedItem.ID] = &item
	}

	// Create a new request
	req, _ := http.NewRequest("GET", "/items", nil)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	GetItemsHandler(&items, &itemMU)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	var itemRes []ItemResponse
	json.NewDecoder(rr.Body).Decode(&itemRes)

	sort.Slice(itemRes, func(i, j int) bool {
		return itemRes[i].ID < itemRes[j].ID
	})

	if !reflect.DeepEqual(itemRes, expectedItemRes) {
		t.Errorf("Expected item responses %+v, but got %+v", expectedItemRes, itemRes)
	}
}

func TestDetermineWinner(t *testing.T) {
	caseL := []struct {
		item        *Item
		winnerIndex int
		highestBid  int
	}{
		{
			item: &Item{
				ID:   1,
				Name: "Test Item - 1",
				Bids: []Bid{
					{BidderName: "Sasha", InitialBid: dollarsToCents(50.00), MaxBid: dollarsToCents(80.00), BidIncrement: dollarsToCents(3.00)},
					{BidderName: "John", InitialBid: dollarsToCents(60.00), MaxBid: dollarsToCents(82.00), BidIncrement: dollarsToCents(2.00)},
					{BidderName: "Pat", InitialBid: dollarsToCents(55.00), MaxBid: dollarsToCents(85.00), BidIncrement: dollarsToCents(5.00)},
				},
			},
			winnerIndex: 2,
			highestBid:  dollarsToCents(85.00),
		},
		{
			item: &Item{
				ID:   2,
				Name: "Test Item - 2",
				Bids: []Bid{
					{BidderName: "Riley", InitialBid: dollarsToCents(700.00), MaxBid: dollarsToCents(725.00), BidIncrement: dollarsToCents(2.00)},
					{BidderName: "Morgan", InitialBid: dollarsToCents(599.00), MaxBid: dollarsToCents(725.00), BidIncrement: dollarsToCents(15.00)},
					{BidderName: "Charlie", InitialBid: dollarsToCents(625.00), MaxBid: dollarsToCents(725.00), BidIncrement: dollarsToCents(8.00)},
				},
			},
			winnerIndex: 0,
			highestBid:  dollarsToCents(722.00),
		},
		{
			item: &Item{
				ID:   3,
				Name: "Test Item - 3",
				Bids: []Bid{
					{BidderName: "Alex", InitialBid: dollarsToCents(2500.00), MaxBid: dollarsToCents(3000.00), BidIncrement: dollarsToCents(500.00)},
					{BidderName: "Jesse", InitialBid: dollarsToCents(2800.00), MaxBid: dollarsToCents(3100.00), BidIncrement: dollarsToCents(201.00)},
					{BidderName: "Drew", InitialBid: dollarsToCents(2501.00), MaxBid: dollarsToCents(3200.00), BidIncrement: dollarsToCents(247.00)},
				},
			},
			winnerIndex: 1,
			highestBid:  dollarsToCents(3001.00),
		},
		{
			item: &Item{
				ID:   4,
				Name: "Test Item - 4",
				Bids: []Bid{
					{BidderName: "Alex", InitialBid: dollarsToCents(2500.00), MaxBid: dollarsToCents(3000.00), BidIncrement: dollarsToCents(500.00)},
				},
			},
			winnerIndex: 0,
			highestBid:  dollarsToCents(2500.00),
		},
		{
			item: &Item{
				ID:   5,
				Name: "Test Item - 5",
				Bids: []Bid{
					{BidderName: "1", InitialBid: dollarsToCents(50.00), MaxBid: dollarsToCents(80.00), BidIncrement: dollarsToCents(3.00)},
					{BidderName: "2", InitialBid: dollarsToCents(50.00), MaxBid: dollarsToCents(80.00), BidIncrement: dollarsToCents(3.00)},
					{BidderName: "3", InitialBid: dollarsToCents(50.00), MaxBid: dollarsToCents(80.00), BidIncrement: dollarsToCents(3.00)},
					{BidderName: "4", InitialBid: dollarsToCents(50.00), MaxBid: dollarsToCents(80.00), BidIncrement: dollarsToCents(3.00)},
					{BidderName: "5", InitialBid: dollarsToCents(50.00), MaxBid: dollarsToCents(80.00), BidIncrement: dollarsToCents(3.00)},
					{BidderName: "6", InitialBid: dollarsToCents(50.00), MaxBid: dollarsToCents(80.00), BidIncrement: dollarsToCents(3.00)},
					{BidderName: "7", InitialBid: dollarsToCents(50.00), MaxBid: dollarsToCents(80.00), BidIncrement: dollarsToCents(3.00)},
					{BidderName: "8", InitialBid: dollarsToCents(50.00), MaxBid: dollarsToCents(80.00), BidIncrement: dollarsToCents(3.00)},
					{BidderName: "9", InitialBid: dollarsToCents(50.00), MaxBid: dollarsToCents(80.00), BidIncrement: dollarsToCents(3.00)},
					{BidderName: "10", InitialBid: dollarsToCents(50.00), MaxBid: dollarsToCents(80.00), BidIncrement: dollarsToCents(3.00)},
				},
			},
			winnerIndex: 0,
			highestBid:  dollarsToCents(80.00),
		},
	}

	for _, c := range caseL {
		winnerIndex, highestBid := c.item.determineWinner()

		if winnerIndex != c.winnerIndex {
			t.Errorf("Expected winner index %d, but got %d", c.winnerIndex, winnerIndex)
		}

		if highestBid != c.highestBid {
			t.Errorf("Expected highest bid %d, but got %d", c.highestBid, highestBid)
		}
	}
}
func TestGetItemWinner(t *testing.T) {
	items := make(map[int]*Item)
	itemIDCounter := 1

	// Create an item with bids
	item := &Item{
		ID:   itemIDCounter,
		Name: "Test Item",
		Bids: []Bid{
			{BidderName: "Sasha", InitialBid: dollarsToCents(50.00), MaxBid: dollarsToCents(80.00), BidIncrement: dollarsToCents(3.00)},
			{BidderName: "John", InitialBid: dollarsToCents(60.00), MaxBid: dollarsToCents(82.00), BidIncrement: dollarsToCents(2.00)},
			{BidderName: "Pat", InitialBid: dollarsToCents(55.00), MaxBid: dollarsToCents(85.00), BidIncrement: dollarsToCents(5.00)},
		},
	}
	items[itemIDCounter] = item

	// Create a new request
	req, _ := http.NewRequest("GET", "/items/{id}/winner", nil)
	req.SetPathValue("id", "1")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	GetItemWinner(&items)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	var bidRes Bid
	json.NewDecoder(rr.Body).Decode(&bidRes)

	expectedBidRes := Bid{
		BidderName:   "Pat",
		InitialBid:   dollarsToCents(55.00),
		MaxBid:       dollarsToCents(85.00),
		BidIncrement: dollarsToCents(5.00),
		CurrentBid:   dollarsToCents(85.00),
	}
	if !reflect.DeepEqual(bidRes, expectedBidRes) {
		t.Errorf("Expected bid response %+v, but got %+v", expectedBidRes, bidRes)
	}
}

func TestGetItemWinnerInvalidID(t *testing.T) {
	items := make(map[int]*Item)

	// Create a new request with an invalid ID
	req, _ := http.NewRequest("GET", "/items/invalid/winner", nil)

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	GetItemWinner(&items)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGetItemWinnerItemNotFound(t *testing.T) {
	items := make(map[int]*Item)

	// Create a new request with a valid ID but no item
	req, _ := http.NewRequest("GET", "/items/{id}/winner", nil)
	req.SetPathValue("id", "1")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	GetItemWinner(&items)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, rr.Code)
	}
}

func TestGetItemWinnerNoWinnerFound(t *testing.T) {
	items := make(map[int]*Item)
	itemIDCounter := 1

	// Create an item with no bids
	item := &Item{
		ID:   itemIDCounter,
		Name: "Test Item",
		Bids: []Bid{},
	}
	items[itemIDCounter] = item

	// Create a new request
	req, _ := http.NewRequest("GET", "/items/{id}/winner", nil)
	req.SetPathValue("id", "1")

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function
	GetItemWinner(&items)(rr, req)

	// Check the response status code
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, rr.Code)
	}
}
