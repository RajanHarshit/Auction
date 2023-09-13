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
	CurrentPrice float64  `json:"currentPrice"`
	EndTime	    time.Time `json:"endTime"`
}

// Bid represent auction process
type Bid struct {
	BidID     string    `json:"bidId"`
	AdspaceID string    `json:"adspaceID"`
	BidderID  string    `json:"bidderID"`
	Amount    float64   `json:"amount"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Duration  time.Duration `json:"duration"`
}

// Bidder represents identity of person
type Bidder struct {
	BidderID   string `json:"bidderId"`
	Name 	   string `json:"name"`
}

type CurrentBidState struct {
	BidderID   string `json:"bidder"`
	Amount     float64 `json:"amount"`
	DateTime   time.Time `json:"dateTime"`
}
