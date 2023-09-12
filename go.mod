module auction.example

replace (
	auction.example/supply_side/handlers => ../go/pkg/mod/auction.example/supply_side/handlers
	auction.example/supply_side/models => ../go/pkg/mod/auction.example/supply_side/models
)

go 1.20

require (
	auction.example/supply_side/models v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.7.1
)
