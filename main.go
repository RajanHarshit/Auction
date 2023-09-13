package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	models "auction.example/auction/models"
	dbconnection "auction.example/dbConnection"
)

var db *sql.DB

// CreateAdSpace creates a new ad space
func CreateAdSpace(w http.ResponseWriter, r *http.Request) {
	var adspace models.Adspace
	_ = json.NewDecoder(r.Body).Decode(&adspace)
	adspace.EndTime = time.Now().UTC().Add(5 * time.Minute)
	adspace.CurrentPrice = adspace.BasePrice

	// table creation
	var query string
	query = "CREATE TABLE IF NOT EXISTS AdSpace ( AdspaceID VARCHAR(255) NOT NULL PRIMARY kEY, Name VARCHAR(255), Description TEXT, BasePrice FLOAT DEFAULT 10, CurrentPrice FLOAT, EndTime TIME)"
	create, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer create.Close()

	// insertion
	query = `INSERT INTO AdSpace (AdspaceID,Name,Description,BasePrice, CurrentPrice, EndTime) VALUES (?, ?, ?, ?, ?, ?)`
	insertResult, err := db.Exec(query, adspace.AdspaceID, adspace.Name, adspace.Description, adspace.BasePrice, adspace.CurrentPrice, adspace.EndTime)

	if err != nil {
		log.Fatalf("impossible insert adspace: %s", err)
	}
	_, err = insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
	}

	// Return a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(adspace)
}

// createBidder creates a bidder name with specific id
func CreateBidder(w http.ResponseWriter, r *http.Request) {
	var query string
	var bidder models.Bidder
	_ = json.NewDecoder(r.Body).Decode(&bidder)

	// table creation
	query = "CREATE TABLE IF NOT EXISTS bidder ( BidderID VARCHAR(255), Name VARCHAR(255))"
	create, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer create.Close()

	// insertion
	query = `INSERT INTO bidder (BidderID, Name) VALUES (?, ?)`
	insertResult, err := db.Exec(query, bidder.BidderID, bidder.Name)
	if err != nil {
		log.Fatalf("impossible insert teacher: %s", err)
	}
	_, err = insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
	}

	// Return a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bidder)
}

// createBid creates a bidder name with specific id
func CreateBid(w http.ResponseWriter, r *http.Request) {
	var query string
	var bid models.Bid
	var adspace models.Adspace
	var bidder models.Bidder
	_ = json.NewDecoder(r.Body).Decode(&bid)
	bid.StartTime = time.Now().UTC()
	bid.EndTime = bid.StartTime.Add(bid.Duration * time.Minute)

	// check adspace available or not
	query = `SELECT AdspaceID, CurrentPrice FROM AdSpace WHERE AdspaceID = ?`
	row, err := db.Query(query, bid.AdspaceID)
	if err != nil {
		log.Fatalf("Unable to Execute query")
		return
	}
	defer row.Close()
	for row.Next() {
		if err := row.Scan(&adspace.AdspaceID, &adspace.CurrentPrice); err != nil {
			log.Fatalf("Error %s when scaning rows \n", err)
			return
		}
	}
	err = row.Err()
	if err != nil {
		log.Fatalf("Error in rows %s\n", err)
		return
	}

	// Check bidder present or not
	query = `SELECT BidderID FROM bidder WHERE BidderID = ?`
	bidderRow := db.QueryRow(query, bid.BidderID)
	if err := bidderRow.Scan(&bidder.BidderID); err != nil {
		log.Fatalf("Error %s when scaning rows \n", err)
		return
	}
	err = bidderRow.Err()
	if err != nil {
		log.Fatalf("Error in rows %s\n", err)
		return
	}

	go MonitorAuctionProcess(&bid)
	currentAuctionState, err := StartAuction(&bid, &adspace)
	if err != nil {
		log.Fatalf("%s", err)
		return
	}

	// table creation
	query = "CREATE TABLE IF NOT EXISTS bid ( BidID VARCHAR(255), AdspaceID VARCHAR(255), BidderID VARCHAR(255), Amount DECIMAL(10, 2), StartTime TIME, EndTime TIME)"
	create, err := db.Query(query)
	if err != nil {
		log.Fatalf("Error in query execution %s\n", err)
		return
	}
	defer create.Close()

	// insertion
	query = `INSERT INTO bid (BidID, AdspaceID, BidderID, Amount, StartTime, EndTime) VALUES (?, ?, ?, ?, ?, ?)`
	insertResult, err := db.Exec(query, bid.BidID, bid.AdspaceID, bid.BidderID, bid.Amount, bid.StartTime, bid.EndTime)
	if err != nil {
		log.Fatalf("impossible insert bidder info: %s", err)
		return
	}
	_, err = insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
		return
	}

	// Return a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(currentAuctionState)
}

func StartAuction(auction *models.Bid, adspace *models.Adspace) (*models.CurrentBidState, error) {
	currentBidState := &models.CurrentBidState{}

	if HasAuctionEnded(auction) {
		return nil, fmt.Errorf("error: auction has been done")
	}

	if auction.Amount <= adspace.CurrentPrice {

		return nil, fmt.Errorf("amount is less than current Price")
	}
	currentBidState.BidderID = auction.BidderID
	currentBidState.Amount = auction.Amount
	currentBidState.DateTime = time.Now().UTC()
	adspace.CurrentPrice = auction.Amount
	return currentBidState, nil
}
func ListAdspace(w http.ResponseWriter, r *http.Request) {
	var adspace models.Adspace
	var arr []interface{}

	defer db.Close()

	// check adspace available or not
	query := `SELECT AdspaceID, Name, Description, BasePrice FROM AdSpace`
	row, err := db.Query(query)
	if err != nil {
		log.Fatalf("Unable to Execute query")
		return
	}
	defer row.Close()
	for row.Next() {
		if err := row.Scan(&adspace.AdspaceID, &adspace.Name, &adspace.Description, &adspace.BasePrice); err != nil {
			log.Fatalf("Unable to scan row")
			return
		}
		arr = append(arr, adspace)
	}
	err = row.Err()
	if err != nil {
		log.Fatalf("Unable to execute row %s", err)
		return
	}

	// Return a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(arr)
}

func ListBidder(w http.ResponseWriter, r *http.Request) {
	var bidder models.Bidder
	var arr []interface{}

	defer db.Close()

	// check adspace available or not
	query := `SELECT * FROM bidder`
	row, err := db.Query(query)
	if err != nil {
		log.Fatalf("Unable to Execute query")
		return
	}
	defer row.Close()
	for row.Next() {
		if err := row.Scan(&bidder.BidderID, &bidder.Name); err != nil {
			log.Fatalf("Unable to scan row")
			return
		}
		arr = append(arr, bidder)
	}
	err = row.Err()
	if err != nil {
		log.Fatalf("Unable to execute row %s", err)
		return
	}

	// Return a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(arr)
}

func HasAuctionEnded(auc *models.Bid) bool {
	currentTime := time.Now().UTC()
	return currentTime.After(auc.EndTime)
}

func MonitorAuctionProcess(auction *models.Bid) {
	if !HasAuctionEnded(auction) {
		time.Sleep(time.Minute)
	}
}

func main() {

	var err error

	// DB connection
	db, err = dbconnection.Connection()
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/createAdSpace", CreateAdSpace)
	mux.HandleFunc("/listAdspace", ListAdspace)
	mux.HandleFunc("/createBidder", CreateBidder)
	mux.HandleFunc("/listBidder", ListBidder)
	mux.HandleFunc("/createBid", CreateBid)
	http.ListenAndServe(":8080", mux)

	defer db.Close()
}
