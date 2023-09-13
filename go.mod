module auction.example

replace (
	auction.example/auction/models => ../go/pkg/mod/auction.example/auction/models
	auction.example/dbConnection => ../go/pkg/mod/auction.example/dbConnection
	auction.example/utils => ../go/pkg/mod/auction.example/utils
)

go 1.20

require (
	auction.example/auction/models v0.0.0-00010101000000-000000000000
	auction.example/dbConnection v0.0.0-00010101000000-000000000000
)

require (
	auction.example/utils v0.0.0-00010101000000-000000000000 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
)
