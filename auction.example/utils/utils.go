package utils

import "fmt"

const (
	username = "root"
	password = "root_password"
	hostname = "127.0.0.1:3306"
	dbname   = "auction_db"
)

func URL(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}
