package dbconnection

import (
	"database/sql"
	"errors"

	"auction.example/utils"
	_ "github.com/go-sql-driver/mysql"
)

const dbname = "auction_db"

func Connection() (*sql.DB, error) {
	//db connection save the ad space to the database
	db, err := sql.Open("mysql", utils.URL(dbname))
	if err != nil {
		return nil, errors.New("error when opening database")
	}
	return db, nil
}
