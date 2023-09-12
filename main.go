package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	models "auction.example/supply_side/models"
)

const (
	username = "root"
	password = "root_password"
	hostname = "127.0.0.1:3306"
	dbname   = "auction_db"
)

// CreateAdSpace creates a new ad space
func CreateAdSpace(w http.ResponseWriter, r *http.Request) {
	var adspace models.Adspace
	_ = json.NewDecoder(r.Body).Decode(&adspace)
	adspace.EndTime = time.Now().UTC().Add(15 * time.Minute)

	//db connection save the ad space to the database
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}

	defer db.Close()

	// table creation
	var query string
	query = "CREATE TABLE IF NOT EXISTS AdSpace ( AdspaceID VARCHAR(255) NOT NULL PRIMARY kEY, Name VARCHAR(255), Description TEXT, BasePrice INT DEFAULT 10, AuctionID TEXT, EndTime TIME)"
	create, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer create.Close()

	// insertion
	query = `INSERT INTO AdSpace (AdspaceID,Name,Description,BasePrice,EndTime) VALUES (?, ?, ?, ?, ?)`
	insertResult, err := db.Exec(query, adspace.AdspaceID, adspace.Name, adspace.Description, adspace.BasePrice, adspace.EndTime)

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

// createAuction create a auction service
func CreateAuction(w http.ResponseWriter, r *http.Request) {
	var query string
	var aucSpace models.Auction
	var adspace models.Adspace

	_ = json.NewDecoder(r.Body).Decode(&aucSpace)
	aucSpace.StartTime = time.Now().UTC()
	aucSpace.EndTime = time.Now().UTC()

	//db connection save the ad space to the database
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}

	defer db.Close()

	// check ad space for auction
	query = `SELECT AdspaceID FROM AdSpace WHERE AdspaceID = ?`
	row := db.QueryRow(query, aucSpace.AdspaceID)
	if err := row.Scan(&adspace.AdspaceID); err != nil {
		panic(err)
	}
	err = row.Err()
	if err != nil {
		log.Fatal(err)
	}

	// table creation
	/*If adspace available then we'll perform auction*/
	query = "CREATE TABLE IF NOT EXISTS auction (ID VARCHAR(255) NOT NULL PRIMARY KEY, StartTime TIME, EndTime TIME, Status VARCHAR(255), AdspaceID VARCHAR(255))"
	create, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer create.Close()

	/*insert record to auction table*/
	query = `INSERT INTO auction (ID, StartTime, EndTime, Status, AdspaceID) VALUES (?, ?, ?, ? , ?)`
	insertResult, err := db.Exec(query, aucSpace.ID, aucSpace.StartTime, aucSpace.EndTime, aucSpace.Status, aucSpace.AdspaceID)
	if err != nil {
		log.Fatalf("impossible insert auction: %s", err)
	}
	_, err = insertResult.LastInsertId()
	if err != nil {
		log.Fatalf("impossible to retrieve last inserted id: %s", err)
	}

	/*Update adspace table with auctionID*/
	query = `UPDATE AdSpace SET AuctionID = ? WHERE AdspaceID = ?`
	updateResult, err := db.Exec(query, aucSpace.ID, adspace.AdspaceID)
	if err != nil {
		log.Fatalf("impossible to update auction: %s", err)
	}
	_, err = updateResult.RowsAffected()
	if err != nil {
		log.Fatalf("Error : unable to update record : %s", err)
	}

	// Return a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(aucSpace)
}

// createBidder creates a bidder name with specific id
func CreateBidder(w http.ResponseWriter, r *http.Request) {
	var query string
	var bidder models.Bidder
	_ = json.NewDecoder(r.Body).Decode(&bidder)

	//db connection save the ad space to the database
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}

	defer db.Close()

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
	bid.Timestamp = time.Now().UTC()

	//db connection save the ad space to the database
	db, err := sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}

	defer db.Close()

	// check adspace available or not
	query = `SELECT AdspaceID FROM AdSpace WHERE AdspaceID = ?`
	row := db.QueryRow(query, bid.AdspaceID)
	if err := row.Scan(&adspace.AdspaceID); err != nil {
		panic(err)
	}
	err = row.Err()
	if err != nil {
		log.Fatal(err)
	}

	// Check bidder present or not
	query = `SELECT BidderID FROM bidder WHERE BidderID = ?`
	bidderRow := db.QueryRow(query, bid.BidderID)
	if err := bidderRow.Scan(&bidder.BidderID); err != nil {
		panic(err)
	}
	err = bidderRow.Err()
	if err != nil {
		log.Fatal(err)
	}

	// table creation
	query = "CREATE TABLE IF NOT EXISTS bid ( BidID VARCHAR(255), AdspaceID VARCHAR(255), BidderID VARCHAR(255), Amount DECIMAL(10, 2), Timestamp TIME)"
	create, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer create.Close()

	// insertion
	query = `INSERT INTO bid (BidID, AdspaceID, BidderID, Amount, Timestamp) VALUES (?, ?, ?, ?, ?)`
	insertResult, err := db.Exec(query, bid.BidID, bid.AdspaceID, bid.BidderID, bid.Amount, bid.Timestamp)
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
	json.NewEncoder(w).Encode(bid)
}

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/createAdSpace", CreateAdSpace)
	mux.HandleFunc("/createAuction", CreateAuction)
	mux.HandleFunc("/createBidder", CreateBidder)
	mux.HandleFunc("/createBid", CreateBid)
	http.ListenAndServe(":8080", mux)
}
