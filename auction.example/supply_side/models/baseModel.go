package models

import (
	"time"
)

// Adspace represents an ad space available for auction
type Adspace struct {
	AdspaceID   string    `json:"adspaceId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BasePrice   float64   `json:"basePrice"`
	AuctionID   string    `json:"auctionID"`
	EndTime	    time.Time `json:"endTime"`
}

// Auction represent an auction for an ad space
type Auction struct {
	ID        string       `json:"id"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Status    string    `json:"status"`
	AdspaceID string `json:"adspaceId"`
}

type Bid struct {
	BidID     string    `json:"bidId"`
	AdspaceID string    `json:"adspaceID"`
	BidderID  string    `json:"bidderID"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type Bidder struct {
	BidderID   string `json:"bidderId"`
	Name 	   string `json:"name"`
}
