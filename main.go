package main

import (
	"fmt"
	"github.com/eager7/ethereum_contracts/database"
)

func main() {
	db, err := database.Initialize("127.0.0.1:3306", "root", "", "eth_database", 2000, 1000)
	if err != nil {
		panic(err)
	}
	contracts, err := database.SearchContracts(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("search contracts len:", len(contracts))
}
